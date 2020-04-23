package core

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	mssqlx "github.com/linxGnu/mssqlx"
)

// InitDBForTestUnsafe init and inject database connection for test (unsafe for concurrent use)
func InitDBForTestUnsafe() (*mssqlx.DBs, error) {
	db, errs := mssqlx.ConnectMasterSlaves(serverConf.Database.Type, []string{testDSN}, []string{testDSN})
	return db, errs[0]
}

// InitDBForProduction init db for production
func InitDBForProduction(c *MysqlConnConfig) (*mssqlx.DBs, error) {
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
