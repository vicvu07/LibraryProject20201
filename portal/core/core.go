package core

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-sql-driver/mysql"
	"github.com/gocql/gocql"
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
	testDSN = "root:123@/test?charset=utf8&collation=utf8_general_ci&parseTime=true&loc=Asia%2FHo_Chi_Minh"
)

var (
	ServiceName   = "portal"
	ErrBadRequest = fmt.Errorf("Bad request")

	ErrExtTermChanCapInvalid = fmt.Errorf("Term chan capacity is invalid")

	// ErrDBObjNull indicate DB Object is nil
	ErrDBObjNull = fmt.Errorf("DB Object is nil")
)

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
	cache          atomic.Value
	grpcServerConf atomic.Value
	httpServerConf atomic.Value // *configs.HTTPClientConf
	loggerLock     = &sync.Mutex{}
	lg             *model.LogFormat
	numShard       uint64
)

type CassandraCache struct {
	Cluster  string `json:"Cluster"`
	Keyspace string `json:"Keyspace"`
}

// SetNumShards set number of shards.
func SetNumShards(v int) {
	atomic.StoreUint64(&numShard, uint64(v))
}

// GetNumShards get number of shards.
func GetNumShards() uint64 {
	return atomic.LoadUint64(&numShard)
}

func InitCore(shardNumber int) {
	SetNumShards(shardNumber)
	lg = logger.MustGet(ServiceName)

	config.Store(&Config{
		Database: configs.MysqlConnConfig{
			Type:    "mysql",
			Masters: []string{testDSN},
			Slaves:  []string{testDSN},
		},
		Cache: CassandraCache{
			Cluster:  "127.0.0.1",
			Keyspace: "portal",
		},
		WebServer: WebServer{
			BodyLimit: "1M",
			Secure: Secure{
				SecureCookie: SecureCookie{
					CookieName:     "auth",
					ContextKey:     "lg_uid",
					MaxAge:         0,
					ExpireInMinute: 1440,
					HashKey:        "Wpdhgkwkpwngapoge93nx9sj2lsigwnx9529xn#px02naigm2-1$93*7nwlwsddf",
					BlockKey:       "DSgjo23058nxg84ns592*383nboiwkg+",
				},
				JWT: JWTConfig{
					ContextKey:     "lg_ctx_key",
					SecretKey:      "Wsjfi^sgjlkajskgwmsmvlk!utwEpc03q<qoP6[%4",
					ExpireInMinute: 1440,
				},
				SipHashSum0: 947295729583939162,
				SipHashSum1: 323869573058327753,
			},
		},
	})

	httpServerConf.Store(&configs.HTTPServerConf{
		Port: 11001,
	})

	grpcServerConf.Store(&configs.GRPCServerConfig{
		PublicIP: "",
		Port:     11000,
	})

	mysqlLibs.RegisterDial()
	constructCassandraSession()
}

func constructCassandraSession() {
	cf := GetConfig()

	cluster := gocql.NewCluster(cf.Cache.Cluster)
	cluster.Keyspace = cf.Cache.Keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	cache.Store(session)
}

// GetDB get global database conn
func GetCacheSession() *gocql.Session {
	if _cache, stored := cache.Load().(*gocql.Session); stored {
		return _cache
	}
	return nil
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
	// WebServer configuration
	WebServer WebServer      `json:"WebServer"` // WebServer configuration
	Cache     CassandraCache `json:"Cache"`
}

// ------------------------------------------------------------- Web Server -------------------------------
// WebServer hold configurations for WebServer
type WebServer struct {
	// BodyLimit The body limit is determined based on both Content-Length request header and actual content read, which makes it super secure.
	// Limit can be specified as 4x or 4xB, where x is one of the multiple from K, M, G, T or P. Example: 2M = 2 Megabyte
	BodyLimit string
	// Secure configuration
	Secure Secure
}

// JWTConfig configuration for jwt token within web server
type JWTConfig struct {
	// ContextKey to get JWT token from context
	ContextKey string
	// SecretKey to generate JWT Token
	SecretKey string
	// ExpireInMinute jwt token will expire after minutes
	ExpireInMinute int64
}

// SecureCookie secure cookie configuration
type SecureCookie struct {
	// CookieName name of secure cookie
	CookieName string
	// ContextKey to get SecureCookie from context
	ContextKey string
	// MaxAge of cookie
	MaxAge int
	// ExpireInMinute
	ExpireInMinute int64
	// HashKey 64 character
	HashKey string
	// BlockKey 32 character
	BlockKey string
}

// Secure config
type Secure struct {
	// JWT for web application/mobile
	JWT JWTConfig

	// SecureCookie secure cookie configuration
	SecureCookie SecureCookie

	// SipHashSum0
	SipHashSum0 uint64
	// SipHashSum1
	SipHashSum1 uint64
}

// Sign jwt token
func (c *JWTConfig) Sign(claim *model.Claim) (string, error) {
	if claim == nil {
		return "", fmt.Errorf("Claim is nil")
	}

	// modify claim for expire at
	claim.ExpiresAt = time.Now().Add(time.Minute * time.Duration(c.ExpireInMinute)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// generate encoded token and send it as response.
	return token.SignedString([]byte(c.SecretKey))
}

// SipHash do sip hash 4-8 sum
func (c *Secure) SipHash(payload []byte) uint64 {
	//	return uint64(lib.SipHash48(c.SipHashSum0, c.SipHashSum1, payload))
	return 0
}

//----------------------------------------------------------------------------------- Usermanager config---------------------------------
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

// GetHTTPServerConf get webserver binding configuration.
func GetHTTPServerConf() *configs.HTTPServerConf {
	if v, stored := httpServerConf.Load().(*configs.HTTPServerConf); stored {
		return v
	}
	return nil
}

// SetHTTPServerConf set webserver binding configuration.
func SetHTTPServerConf(v *configs.HTTPServerConf) {
	httpServerConf.Store(v)
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
