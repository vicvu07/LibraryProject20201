package core

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-multierror"
	"github.com/linxGnu/mssqlx"

	"github.com/pinezapple/LibraryProject20201/skeleton/configs"
	etcd "github.com/pinezapple/LibraryProject20201/skeleton/libs/etcd"
	mysqlLibs "github.com/pinezapple/LibraryProject20201/skeleton/libs/mysql"
	"github.com/pinezapple/LibraryProject20201/skeleton/logger"
	"github.com/pinezapple/LibraryProject20201/skeleton/model"
	etcdClient "go.etcd.io/etcd/clientv3"
)

const (
	testDSN     = "root:123@/test?charset=utf8&collation=utf8_general_ci&parseTime=true&loc=Asia%2FHo_Chi_Minh"
	ServiceName = "docmanager"
)

var lg *model.LogFormat

// ------------------------- Etcd Client -------------------------
var eClient atomic.Value

// GetEtcdClient get global etcd client.
func GetEtcdClient() *etcdClient.Client {
	if client, stored := eClient.Load().(*etcdClient.Client); stored {
		return client
	}
	return nil
}

// SetEtcdClient set glocal etcd client.
func SetEtcdClient(e *etcdClient.Client) {
	eClient.Store(e)
}

// ------------------------- Main Config -------------------------
var (
	config         atomic.Value
	db             atomic.Value
	grpcServerConf atomic.Value
	loggerLock     = &sync.Mutex{}
)

func init() {
	lg = logger.MustGet(ServiceName)
	config.Store(&Config{
		Database: configs.MysqlConnConfig{
			Type:    "mysql",
			Masters: []string{testDSN},
			Slaves:  []string{testDSN},
		},
	})

	grpcServerConf.Store(&configs.GRPCServerConfig{
		PublicIP: "",
		Port:     10000,
	})

	mysqlLibs.RegisterDial()
}

func GetLogger() (_lg *model.LogFormat) {
	loggerLock.Lock()
	_lg = lg
	loggerLock.Unlock()
	return
}

// Config configuration of campaign
type Config struct {
	// Database configuration
	Database configs.MysqlConnConfig `json:"Database"`
}

// GetConfig get global config
func GetConfig() *Config {
	if cnf, stored := config.Load().(*Config); stored {
		return cnf
	}
	return nil
}

// GetDB get global database conn
func GetDB() *mssqlx.DBs {
	if _db, stored := db.Load().(*mssqlx.DBs); stored {
		return _db
	}
	return nil
}

// GetGRPCServerConfig get grpc server configuration.
func GetGRPCServerConfig() *configs.GRPCServerConfig {
	if cnf, stored := grpcServerConf.Load().(*configs.GRPCServerConfig); stored {
		return cnf
	}
	return nil
}

// SetGRPCServerConfig set grpc server configuration.
func SetGRPCServerConfig(cf *configs.GRPCServerConfig) {
	grpcServerConf.Store(cf)
}

// InitDBForTest init and inject database connection for test.
func InitDBForTest() {
	conf := GetConfig()
	_db, _ := mssqlx.ConnectMasterSlaves(conf.Database.Type, conf.Database.Masters, conf.Database.Slaves)
	db.Store(_db)
}

// WatchConfig watch configuration.
func WatchConfig(ctx context.Context) (model.Daemon, error) {
	// check etcd client
	cl := GetEtcdClient()
	if cl == nil {
		//lg.Err("etcd client not init")
		logger.LogInfo(lg, "etcd client not init")
		panic("etcd client fault")
	}

	return etcd.Watch(ctx, cl, ServiceName, GetConfigFromEtcd)
}

// GetConfigFromEtcd get and update global configuration from etcd
func GetConfigFromEtcd() (err error) {
	// check etcd client
	cl := GetEtcdClient()
	if cl == nil {
		//lg.Fatal("etcd client not init")
		logger.LogInfo(lg, "etcd client not init")
		panic("etcd client fault")
	}

	fmt.Println(cl)
	return etcd.Get(cl, ServiceName, &Config{}, func(expect interface{}) (err error) {
		if expect == nil {
			return
		}

		obj, ok := expect.(*Config)
		if !ok {
			return
		}

		fmt.Println(obj)
		defer func() {
			if err == nil {
				config.Store(obj)
			}
		}()

		// try to initialize db
		initializeDB := func(c *configs.MysqlConnConfig) (*mssqlx.DBs, error) {
			if c == nil {
				return nil, nil
			}
			if len(c.Masters) == 0 || len(c.Slaves) == 0 {
				return nil, fmt.Errorf("Masters and Slaves must not be empty")
			}

			dbs, errs := mssqlx.ConnectMasterSlaves(c.Type, c.Masters, c.Slaves, c.IsWsrep)

			nMasters := len(c.Masters)
			var err error

			masterOK, slaveOK := 0, 0
			for i := range errs {
				if errs[i] == nil {
					if i < nMasters {
						masterOK++
					} else {
						slaveOK++
					}
				} else {
					err = multierror.Append(err, errs[i])
				}
			}

			if masterOK == 0 || slaveOK == 0 {
				if dbs != nil {
					dbs.Destroy()
				}
				return nil, err
			}

			return dbs, nil
		}

		// try database construct
		if err = obj.Database.Construct(mysql.RegisterTLSConfig); err != nil {
			return
		}

		// validate database configuration
		if _config, _db := GetConfig(), GetDB(); _config != nil && obj.Database.Equal(&_config.Database) && _db != nil {
			_db.SetMaxIdleConns(obj.Database.MaxIdleConn)
			_db.SetMaxOpenConns(obj.Database.MaxOpenConn)
		} else {
			// try to connect first
			newDB, e := initializeDB(&obj.Database)
			if e != nil {
				err = e
				return
			}
			newDB.SetMaxIdleConns(obj.Database.MaxIdleConn)
			newDB.SetMaxOpenConns(obj.Database.MaxOpenConn)

			// destroy old db
			if _db != nil {
				_db.Destroy()
			}
			db.Store(newDB)
		}
		return
	})
}
