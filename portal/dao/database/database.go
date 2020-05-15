package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pinezapple/LibraryProject20201/portal/core"
	"github.com/pinezapple/LibraryProject20201/skeleton/model"

	"github.com/linxGnu/mssqlx"
	"golang.org/x/crypto/bcrypt"
)

const (
	sqlUserSelect                   = "SELECT * FROM user WHERE id = ?"
	sqlUserSelectByUsername         = "SELECT * FROM user WHERE username = ?"
	sqlUserSecSelectByUsername      = "SELECT * FROM user_security WHERE username = ?"
	sqlUserInsert                   = "INSERT INTO user (username, status, checksum, name, role, dob, sex, phonenumber) VALUES (?,?,?,?,?,?,?,?)"
	sqlUserInsertSec                = "INSERT INTO user_security (username, password, checksum) VALUES (?,?,?)"
	sqlUserUpdateInfo               = "UPDATE user SET name = ?, role = ?, dob = ?, sex = ?, phonenumber = ?, updated_at = ? WHERE username = ?"
	sqlUserSecUpdateWithPassword    = "UPDATE user_security SET password = ?, checksum = ?, updated_at = ? WHERE username = ?"
	sqlUserSecUpdateWithoutPassword = "UPDATE user_security SET checksum = ?, updated_at = ? WHERE username = ?"
	sqlUserUpdateStatus             = "UPDATE user SET status = ?, updated_at = ? WHERE id = ?"
	sqlUserUpdatePasswordCounter    = "UPDATE user SET pwd_counter = ?, updated_at = ? WHERE id = ?"
)

// IUserDAO userDAO interface
type IUserDAO interface {
	// Select select user by id
	Select(ctx context.Context, db *mssqlx.DBs, id uint64) (result *model.User, err error)
	// SelectByUsername select user by username
	SelectByUsername(ctx context.Context, db *mssqlx.DBs, username string) (result *model.User, err error)
	// SelectSecByUsername select user sec by username
	SelectSecByUsername(ctx context.Context, db *mssqlx.DBs, username string) (result *model.UserSecurity, err error)
	// Create create new user
	Create(ctx context.Context, db *mssqlx.DBs, k0, k1 uint64, username, password string, name, role, dob, sex, phonenumber string) (result *model.User, err error)
	// Update update user
	Update(ctx context.Context, db *mssqlx.DBs, k0, k1 uint64, id uint64, password string, name, role, dob, sex, phonenumber string) (err error)
	// UpdateStatus update user status
	UpdateStatus(ctx context.Context, db *mssqlx.DBs, id uint64, status byte) (err error)
	// UpdatePasswordCounter update user password counter
	UpdatePasswordCounter(ctx context.Context, db *mssqlx.DBs, id uint64, counter int) (err error)
}

type userDAO struct{}

// Select select user by id
func (c *userDAO) Select(ctx context.Context, db *mssqlx.DBs, id uint64) (result *model.User, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	result = &model.User{}
	if err = db.GetContext(ctx, result, sqlUserSelect, id); err == sql.ErrNoRows {
		result, err = nil, nil
		return
	}

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

// Create create user
func (c *userDAO) Create(ctx context.Context, db *mssqlx.DBs, k0, k1 uint64, username, password string, name, role, dob, sex, phonenumber string) (result *model.User, err error) {
	//	fmt.Println("init account root")
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		//		fmt.Println("init account root 1")
		return
	}

	user := &model.User{
		Username: username,
		// TODO: using const later
		Name:        name,
		Role:        role,
		Dob:         dob,
		Sex:         sex,
		PhoneNumber: phonenumber,
		Status:      0,
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
		//		fmt.Println("init account root 2")
		return
	}

	if err = model.ExecTransaction(ctx, tx, func(ctx context.Context, tx *sql.Tx) (er error) {
		if _, er = tx.ExecContext(ctx, sqlUserInsert, user.Username, user.Status, user.Checksum, user.Name, user.Role, user.Dob, user.Sex, user.PhoneNumber); er != nil {
			//			fmt.Println("init account root 3")
			return
		}
		if _, er = tx.ExecContext(ctx, sqlUserInsertSec, sec.Username, sec.Password, sec.Checksum); er != nil {
			//fmt.Println(er)
			return
		}

		return
	}); err != nil {
		//fmt.Println("init account root 5")
		return
	}

	// FIX: select on master instead of slave
	result = &model.User{}
	if err = db.GetContextOnMaster(ctx, result, sqlUserSelectByUsername, username); err == sql.ErrNoRows {
		return
	}

	// TODO: using const later
	result.Status = 1
	result.Checksum = result.Sum(k0, k1)

	_, err = db.ExecContext(ctx, "UPDATE user SET checksum = ?, status = ? WHERE id = ?", result.Checksum, result.Status, result.ID)

	return
}

// Update update user
func (c *userDAO) Update(ctx context.Context, db *mssqlx.DBs, k0, k1 uint64, id uint64, password string, name, role, dob, sex, phonenumber string) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	// select user
	user, err := c.Select(ctx, db, id)
	if err != nil {
		return err
	} else if user == nil {
		return fmt.Errorf("User not found")
	}

	// select user sec
	userSec, err := c.SelectSecByUsername(ctx, db, user.Username)
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
			if _, er = tx.ExecContext(ctx, sqlUserUpdateInfo, name, role, dob, sex, phonenumber, now, user.Username); er != nil {
				return
			}
			if _, er = tx.ExecContext(ctx, sqlUserSecUpdateWithPassword, userSec.Password, userSec.Checksum, now, user.Username); er != nil {
				return
			}
		} else {
			if _, er = tx.ExecContext(ctx, sqlUserUpdateInfo, name, role, dob, sex, phonenumber, now, user.Username); er != nil {
				return
			}
			if _, er = tx.ExecContext(ctx, sqlUserSecUpdateWithoutPassword, userSec.Checksum, now, user.Username); er != nil {
				return
			}
		}
		return
	})
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
