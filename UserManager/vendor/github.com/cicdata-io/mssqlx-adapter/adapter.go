// Copyright 2017 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mssqlxadapter

import (
	"errors"
	"runtime"

	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	mssqlx "github.com/linxGnu/mssqlx"
)

const (
	sqlCreateTable       = "CREATE TABLE IF NOT EXISTS casbin_rule(`p_type` VARCHAR(100), `v0` VARCHAR(100), `v1` VARCHAR(100), `v2` VARCHAR(100),`v3` VARCHAR(100),`v4` VARCHAR(100),`v5` VARCHAR(100))"
	sqlInsertPolicy      = "INSERT INTO casbin_rule(`p_type`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`) VALUES (?,?,?,?,?,?,?)"
	sqlDeletePolicy      = "DELETE FROM casbin_rule WHERE p_type=? AND v0=? AND v1=? AND v2=? AND v3=? AND v4=? AND v5=?"
	sqlDeleteAllPolicies = "DELETE FROM casbin_rule"
)

type CasbinRule struct {
	PType string `db:"p_type"`
	V0    string `db:"v0"`
	V1    string `db:"v1"`
	V2    string `db:"v2"`
	V3    string `db:"v3"`
	V4    string `db:"v4"`
	V5    string `db:"v5"`
}

func (c *CasbinRule) TableName() string {
	return "casbin_rule" //as Gorm keeps table names are plural, and we love consistency
}

// Adapter represents the Gorm adapter for policy storage.
type Adapter struct {
	driverName  string
	masterDSNs  []string
	slaveDSNs   []string
	dbSpecified bool
	db          *mssqlx.DBs
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
	a.db.Destroy()
}

// NewAdapter is the constructor for Adapter.
// dbSpecified is an optional bool parameter. The default value is false.
// It's up to whether you have specified an existing DB in dataSourceName.
// If dbSpecified == true, you need to make sure the DB in dataSourceName exists.
// If dbSpecified == false, the adapter will automatically create a DB named "casbin".
func NewAdapter(driverName string, masterDSNs []string, slaveDSNs []string, dbSpecified ...bool) *Adapter {
	a := &Adapter{}
	a.driverName = driverName
	a.masterDSNs = masterDSNs
	a.slaveDSNs = slaveDSNs

	if len(dbSpecified) == 0 {
		a.dbSpecified = false
	} else if len(dbSpecified) == 1 {
		a.dbSpecified = dbSpecified[0]
	} else {
		panic(errors.New("invalid parameter: dbSpecified"))
	}

	// Open the DB, create it if not existed.
	a.open()

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)

	return a
}

func (a *Adapter) createDatabase() []error {

	var err []error
	var db *mssqlx.DBs
	if a.driverName == "postgres" {
		for i := 0; i < len(a.masterDSNs); i++ {
			a.masterDSNs[i] = a.masterDSNs[i] + "dbname=postgres"
		}

		for i := 0; i < len(a.masterDSNs); i++ {
			a.slaveDSNs[i] = a.slaveDSNs[i] + "dbname=postgres"
		}
		db, err = mssqlx.ConnectMasterSlaves(a.driverName, a.masterDSNs, a.slaveDSNs)
	} else {
		db, err = mssqlx.ConnectMasterSlaves(a.driverName, a.masterDSNs, a.slaveDSNs)
	}
	//log.Println("im here", err)
	if err != nil {
		return err
	}
	defer db.Destroy()

	if a.driverName == "postgres" {
		if _, er := db.Exec("CREATE DATABASE casbin"); er != nil {
			err = append(err, er)
			panic(err)
		}
	} else if a.driverName != "sqlite3" {
		if _, er := db.Exec("CREATE DATABASE IF NOT EXISTS casbin"); er != nil {
			err = append(err, er)
			panic(er)
		}
	}
	return err
}

func (a *Adapter) open() {
	var err []error
	var db *mssqlx.DBs
	if a.dbSpecified {
		db, err = mssqlx.ConnectMasterSlaves(a.driverName, a.masterDSNs, a.slaveDSNs)
		for _, e := range err {
			if e != nil {
				panic(e)
			}
		}
	} else {
		if err = a.createDatabase(); err != nil {
			panic(err)
		}
		if a.driverName == "postgres" {
			for i := 0; i < len(a.masterDSNs); i++ {
				a.masterDSNs[i] = a.masterDSNs[i] + "dbname=casbin"
			}

			for i := 0; i < len(a.masterDSNs); i++ {
				a.slaveDSNs[i] = a.slaveDSNs[i] + "dbname=casbin"
			}
			db, err = mssqlx.ConnectMasterSlaves(a.driverName, a.masterDSNs, a.slaveDSNs)
		} else {
			for i := 0; i < len(a.masterDSNs); i++ {
				// a.masterDSNs[i] = a.masterDSNs[i] + "/casbin"
				a.masterDSNs[i] = a.masterDSNs[i]
			}

			for i := 0; i < len(a.masterDSNs); i++ {
				// a.slaveDSNs[i] = a.slaveDSNs[i] + "/casbin"
				a.slaveDSNs[i] = a.slaveDSNs[i]
			}
			db, err = mssqlx.ConnectMasterSlaves(a.driverName, a.masterDSNs, a.slaveDSNs)
		}
	}
	a.db = db
	a.createTable()
}

func (a *Adapter) close() {
	a.db.Destroy()
	a.db = nil
}

func (a *Adapter) createTable() {
	_, err := a.db.Exec(sqlCreateTable)
	if err != nil {
		panic(err)
	}
}

func (a *Adapter) dropTable() {
	_, err := a.db.Exec("DROP TABLE casbin_rule")
	if err != nil {
		panic(err)
	}
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	lineText := line.PType
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}
	persist.LoadPolicyLine(lineText, model)
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	var lines []CasbinRule

	err := a.db.Select(&lines, "SELECT *FROM casbin_rule")
	if err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{}

	line.PType = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) error {

	a.dropTable()
	a.createTable()

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			_, err := a.db.Exec(sqlInsertPolicy, line.PType, line.V0, line.V1, line.V2, line.V3, line.V4, line.V5)
			if err != nil {
				return err
			}
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			_, err := a.db.Exec(sqlInsertPolicy, line.PType, line.V0, line.V1, line.V2, line.V3, line.V4, line.V5)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := a.db.Exec(sqlInsertPolicy, line.PType, line.V0, line.V1, line.V2, line.V3, line.V4, line.V5)
	return err
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	err := rawDelete(a.db, line) //can't use db.Delete as we're not using primary key http://jinzhu.me/gorm/crud.html#delete
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := CasbinRule{}

	line.PType = ptype
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
	}
	err := rawDeleteAll(a.db, line)
	return err
}

func rawDelete(db *mssqlx.DBs, line CasbinRule) error {
	_, err := db.Exec(sqlDeletePolicy, line.PType, line.V0, line.V1, line.V2, line.V3, line.V4, line.V5)
	return err
}

func rawDeleteAll(db *mssqlx.DBs, line CasbinRule) error {
	queryArgs := []interface{}{line.PType}
	//var queryArgs []string
	queryStr := sqlDeleteAllPolicies + " where p_type = ?"
	if line.V0 != "" {
		queryStr += " and v0 = ?"
		queryArgs = append(queryArgs, line.V0)
	}
	if line.V1 != "" {
		queryStr += " and v1 = ?"
		queryArgs = append(queryArgs, line.V1)
	}
	if line.V2 != "" {
		queryStr += " and v2 = ?"
		queryArgs = append(queryArgs, line.V2)
	}
	if line.V3 != "" {
		queryStr += " and v3 = ?"
		queryArgs = append(queryArgs, line.V3)
	}
	if line.V4 != "" {
		queryStr += " and v4 = ?"
		queryArgs = append(queryArgs, line.V4)
	}
	if line.V5 != "" {
		queryStr += " and v5 = ?"
		queryArgs = append(queryArgs, line.V5)
	}
	args := append([]interface{}{}, queryArgs...)
	_, err := db.Exec(queryStr, args...)
	return err
}
