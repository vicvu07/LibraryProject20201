package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"DBproject1/core"
	"DBproject1/model"

	"github.com/linxGnu/mssqlx"
	"golang.org/x/crypto/bcrypt"
)

var blankBytes = []byte{}

const (
	sqlUserSelectByUsername    = "SELECT * FROM user WHERE username = ?"
	sqlUserSecSelectByUsername = "SELECT * FROM user_security WHERE username = ?"

	sqlUserDetailSelectAll                = "SELECT id,name,username FROM user_detail"
	sqlUserDetailSelectByView             = "SELECT * FROM user_detail ORDER BY id DESC LIMIT ?, ?"
	sqlUserDetailSelectByID               = "SELECT * FROM user_detail WHERE id = ?"
	sqlUserDetailSelectUserIDWithUsername = "SELECT id FROM user_detail WHERE username = ?"
	sqlUserDetailInsert                   = "INSERT INTO user_detail(name,dob,sex,position,phonenum,national_id,salary,username) VALUES (?,?,?,?,?,?,?,?)"
	sqlSelectDeparmentForUser             = "SELECT * FROM department_management WHERE status = ? AND user_id = ?"
	sqlInsertDepartmentManagement         = "INSERT INTO department_management(department_id,user_id) VALUES (?,?)"
	sqlSelectDepartmentManagement         = "SELECT * FROM department_management WHERE department_id = ? AND user_id = ?"
	sqlUpdateDepartmentManagementStatus   = "UPDATE department_management SET status = ? WHERE department_id = ? AND user_id = ?"
	sqlDeleteUserFromDepartment           = "DELETE FROM department_management WHERE user_id = ?"
	sqlDeleteUserFromPlan                 = "DELETE FROM plan_management WHERE user_id = ?"
	sqlUserDetailUpdate                   = "UPDATE user_detail SET name = ?, dob = ?, sex = ?, position = ?, phonenum = ?, national_id = ?, salary = ?, updated_at = ? WHERE id = ?"
	sqlUserDetailDelete                   = "DELETE FROM user_detail WHERE id = ?"
	sqlSelectUserDonePlan                 = `SELECT plan.id, plan.name FROM plan, plan_management WHERE plan_management.IndiOrDepart = "I" AND plan_management.foreign_id = ? AND plan_management.plan_id=plan.id AND plan.current_status = "D"`
	sqlSelectUserOnGoingPlan              = `SELECT plan.id, plan.name FROM plan, plan_management WHERE plan_management.IndiOrDepart = "I" AND plan_management.foreign_id = ? AND plan_management.plan_id=plan.id AND plan.current_status = "O"`

	sqlUserInsert                   = "INSERT INTO user (username, status, checksum) VALUES (?,?,?)"
	sqlUserInsertSec                = "INSERT INTO user_security (username, password, checksum) VALUES (?,?,?)"
	sqlUserDelete                   = "DELETE FROM user WHERE username = ?"
	sqlUserDeleteSec                = "DELETE FROM user_security WHERE username = ?"
	sqlUserSecUpdateWithPassword    = "UPDATE user_security SET password = ?, checksum = ?, updated_at = ? WHERE username = ?"
	sqlUserSecUpdateWithoutPassword = "UPDATE user_security SET checksum = ?, updated_at = ? WHERE username = ?"
	sqlUserUpdateStatus             = "UPDATE user SET status = ?, updated_at = ? WHERE id = ?"
	sqlUserUpdatePasswordCounter    = "UPDATE user SET pwd_counter = ?, updated_at = ? WHERE id = ?"
)

// IUserDAO userDAO interface
type IUserDAO interface {
	// Select select user by id
	Select(ctx context.Context, db *mssqlx.DBs, id uint64) (result *model.UserDetail, err error)
	// SelectAll select all user
	SelectAll(ctx context.Context, db *mssqlx.DBs) (result []*model.UserDetailForPlan, err error)
	// SelectByUsername select user by username
	SelectByUsername(ctx context.Context, db *mssqlx.DBs, username string) (result *model.User, err error)
	// SelectSecByUsername select user sec by username
	SelectSecByUsername(ctx context.Context, db *mssqlx.DBs, username string) (result *model.UserSecurity, err error)
	// SelectView view of campaigns, order by id desc with offset and limit
	SelectView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.UserDetail, err error)
	// SelectUserDepartment select user current department
	SelectUserDepartment(ctx context.Context, db *mssqlx.DBs, id uint64) (result *model.DepartmentManagement, err error)
	// Create create new user
	Create(ctx context.Context, db *mssqlx.DBs, k0, k1 uint64, name, dob, sex, position, phonenum, nationalid string, salary uint64, username string, departmentID uint64, password string) (result *model.User, err error)
	// Delete delete user
	Delete(ctx context.Context, db *mssqlx.DBs, username string, id uint64) (err error)
	// UpdatePassword update user password
	UpdatePassword(ctx context.Context, db *mssqlx.DBs, k0, k1 uint64, username string, password string) (err error)
	// UpdateUserDetail update user detail
	UpdateUserDetail(ctx context.Context, db *mssqlx.DBs, id uint64, name, dob, sex, position, phonenum, nationalid string, salary uint64) (err error)
	// UpdateUserDepartmentManagement update user department
	UpdateUserDepartmentManagement(ctx context.Context, db *mssqlx.DBs, id, oldDepartmentID, newDepartmentID uint64) (err error)
	// UpdateStatus update user status
	UpdateStatus(ctx context.Context, db *mssqlx.DBs, id uint64, status byte) (err error)
	// UpdatePasswordCounter update user password counter
	UpdatePasswordCounter(ctx context.Context, db *mssqlx.DBs, id uint64, counter int) (err error)
	// SelectUserDonePlan
	SelectUserDonePlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error)
	// SelectUserOnGoingPlan
	SelectUserOnGoingPlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error)
}

type userDAO struct{}

// Select select user by id
func (c *userDAO) Select(ctx context.Context, db *mssqlx.DBs, id uint64) (result *model.UserDetail, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	result = &model.UserDetail{}
	if err = db.GetContext(ctx, result, sqlUserDetailSelectByID, id); err == sql.ErrNoRows {
		result, err = nil, nil
		return
	}

	return
}

func (c *userDAO) SelectAll(ctx context.Context, db *mssqlx.DBs) (result []*model.UserDetailForPlan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlUserDetailSelectAll)

	return
}

// SelectByUsername select user by username
func (c *userDAO) SelectByUsername(ctx context.Context, db *mssqlx.DBs, username string) (result *model.User, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	result = &model.User{}
	if err = db.GetContext(ctx, result, sqlUserSelectByUsername, username); err == sql.ErrNoRows {
		result, err = nil, nil
		return
	}

	return
}

// SelectSecByUsername select user sec by username
func (c *userDAO) SelectSecByUsername(ctx context.Context, db *mssqlx.DBs, username string) (result *model.UserSecurity, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	result = &model.UserSecurity{}
	if err = db.GetContext(ctx, result, sqlUserSecSelectByUsername, username); err == sql.ErrNoRows {
		result, err = nil, nil
		return
	}

	return
}

// SelectView view of campaigns, order by id desc with offset and limit
func (c *userDAO) SelectView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.UserDetail, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlUserDetailSelectByView, offset, limit)

	return
}

// Create create user
func (c *userDAO) Create(ctx context.Context, db *mssqlx.DBs, k0, k1 uint64, name, dob, sex, position, phonenum, nationalid string, salary uint64, username string, departmentID uint64, password string) (result *model.User, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	user := &model.User{
		Username: username,
		Status:   0,
	}

	p, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	sec := &model.UserSecurity{
		Username: username,
		Password: p,
	}
	sec.Checksum = sec.Sum(k0, k1)

	// prepare transaction
	tx, e := db.Begin()
	if e != nil {
		err = e
		return
	}

	if _, er := db.ExecContext(ctx, sqlUserDetailInsert, name, dob, sex, position, phonenum, nationalid, salary, username); er != nil {
		err = er
		return
	}
	var userID uint64
	if err = db.GetContext(ctx, &userID, sqlUserDetailSelectUserIDWithUsername, username); err == sql.ErrNoRows {
		result, err = nil, fmt.Errorf("Can't Insert User Detail")
		return
	}
	if _, er := db.ExecContext(ctx, sqlInsertDepartmentManagement, departmentID, userID); er != nil {
		err = er
		return
	}

	if err = model.ExecTransaction(ctx, tx, func(ctx context.Context, tx *sql.Tx) (er error) {
		if _, er = tx.ExecContext(ctx, sqlUserInsert, user.Username, user.Status, user.Checksum); er != nil {
			return
		}
		if _, er = tx.ExecContext(ctx, sqlUserInsertSec, sec.Username, sec.Password, sec.Checksum); er != nil {
			return
		}
		return
	}); err != nil {
		return
	}

	if result, e = c.SelectByUsername(ctx, db, username); e != nil {
		err = e
		return
	}
	// TODO: using const later
	result.Status = 1
	result.Checksum = result.Sum(k0, k1)

	_, err = db.ExecContext(ctx, "UPDATE user SET checksum = ?, status = ? WHERE id = ?", result.Checksum, result.Status, result.ID)
	return
}

func (c *userDAO) Delete(ctx context.Context, db *mssqlx.DBs, username string, id uint64) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	_, err = db.Exec(sqlUserDetailDelete, id)
	_, err = db.Exec(sqlDeleteUserFromDepartment, id)
	_, err = db.Exec(sqlDeleteUserFromPlan, id)
	_, err = db.Exec(sqlUserDelete, username)
	_, err = db.Exec(sqlUserDeleteSec, username)
	return
}

// Update update user
func (c *userDAO) UpdatePassword(ctx context.Context, db *mssqlx.DBs, k0, k1 uint64, username string, password string) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	// select user sec
	userSec, err := c.SelectSecByUsername(ctx, db, username)
	if err != nil {
		return err
	}

	// try to generate password
	if len(password) > 0 {
		p, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		userSec.Password = p
	}
	userSec.Checksum = userSec.Sum(k0, k1)

	// prepare transaction
	tx, e := db.Begin()
	if e != nil {
		err = e
		return
	}

	err = model.ExecTransaction(ctx, tx, func(ctx context.Context, tx *sql.Tx) (er error) {
		now := time.Now()
		if len(password) > 0 {
			if _, er = tx.ExecContext(ctx, sqlUserSecUpdateWithPassword, userSec.Password, userSec.Checksum, now, username); err != nil {
				return
			}
		} else {
			if _, er = tx.ExecContext(ctx, sqlUserSecUpdateWithoutPassword, userSec.Checksum, now, username); er != nil {
				return
			}
		}
		return
	})
	return
}

func (c *userDAO) UpdateUserDepartmentManagement(ctx context.Context, db *mssqlx.DBs, id, oldDepartmentID, newDepartmentID uint64) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	if oldDepartmentID == newDepartmentID {
		return
	}

	var res = &model.DepartmentManagement{}
	if erro := db.GetContext(ctx, res, sqlSelectDepartmentManagement, newDepartmentID, id); erro == sql.ErrNoRows {
		if _, er := db.ExecContext(ctx, sqlUpdateDepartmentManagementStatus, 0, oldDepartmentID, id); er != nil {
			err = er
			return
		}
		if _, er := db.ExecContext(ctx, sqlInsertDepartmentManagement, newDepartmentID, id); er != nil {
			err = er
			return
		}
	} else {
		if _, er := db.ExecContext(ctx, sqlUpdateDepartmentManagementStatus, 1, newDepartmentID, id); er != nil {
			err = er
			return
		}
		if _, er := db.ExecContext(ctx, sqlUpdateDepartmentManagementStatus, 0, oldDepartmentID, id); er != nil {
			err = er
			return
		}
	}
	return
}

func (c *userDAO) UpdateUserDetail(ctx context.Context, db *mssqlx.DBs, id uint64, name, dob, sex, position, phonenum, nationalid string, salary uint64) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	if _, er := db.ExecContext(ctx, sqlUserDetailUpdate, name, dob, sex, position, phonenum, nationalid, salary, time.Now(), id); er != nil {
		err = er
		return
	}
	return
}

// UpdateStatus update user status
func (c *userDAO) UpdateStatus(ctx context.Context, db *mssqlx.DBs, id uint64, status byte) (err error) {
	_, err = db.ExecContext(ctx, sqlUserUpdateStatus, status, time.Now(), id)
	return
}

// UpdatePasswordCounter update password counter
func (c *userDAO) UpdatePasswordCounter(ctx context.Context, db *mssqlx.DBs, id uint64, counter int) (err error) {
	_, err = db.ExecContext(ctx, sqlUserUpdatePasswordCounter, counter, time.Now(), id)
	return
}

func (c *userDAO) SelectUserDepartment(ctx context.Context, db *mssqlx.DBs, id uint64) (result *model.DepartmentManagement, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	i := 1
	result = &model.DepartmentManagement{}

	if err = db.GetContext(ctx, result, sqlSelectDeparmentForUser, i, id); err == sql.ErrNoRows {
		result, err = nil, nil
		return
	}

	return
}

func (c *userDAO) SelectUserDonePlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectUserDonePlan, id)
	return
}

func (c *userDAO) SelectUserOnGoingPlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectUserOnGoingPlan, id)
	return
}

var uDAO IUserDAO = &userDAO{}

// GetUserDAO get user dao
func GetUserDAO() IUserDAO {
	return uDAO
}

// SetUserDAO set campaign dao
// NOTE: USE THIS FUNC FOR MOCKING
func SetUserDAO(v IUserDAO) {
	uDAO = v
}
