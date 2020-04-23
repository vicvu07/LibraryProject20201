package dao

import (
	"DBproject1/core"
	"DBproject1/model"
	"context"
	"database/sql"
	"time"

	"github.com/linxGnu/mssqlx"
)

const (
	sqlDepartmentSelectAll              = "SELECT id, name FROM department"
	sqlDepartmentSelectByView           = "SELECT * FROM department ORDER BY id DESC LIMIT ?, ?"
	sqlDepartmentSelectByID             = "SELECT * FROM department WHERE id = ?"
	sqlDepartmentSave                   = "INSERT INTO department(name,description,total_salary) value (?,?,?)"
	sqlDepartmentUpdate                 = "UPDATE department SET name = ?, description = ?, total_salary = ?, updated_at = ? WHERE id = ?"
	sqlDepartmentUpdateWithoutSalary    = "UPDATE department SET name = ?, description = ?, updated_at = ? WHERE id = ?"
	sqlSelectWorkingUserInDepartment    = "SELECT user_detail.id, user_detail.name FROM department_management, user_detail WHERE department_management.department_id = ? AND department_management.status = 1 AND department_management.user_id = user_detail.id"
	sqlSelectUsedToWorkUserInDepartment = "SELECT user_detail.id, user_detail.name FROM department_management, user_detail WHERE department_management.department_id = ? AND department_management.status = 0 AND department_management.user_id = user_detail.id"
	sqlSelectOnGoingPlanInDepartment    = `SELECT plan.id, plan.name, plan.description FROM plan, plan_management WHERE plan_management.IndiOrDepart = "D" AND plan_management.foreign_id = ? AND plan_management.plan_id=plan.id AND plan.current_status = "O"`
	sqlSelectDonePlanInDepartment       = `SELECT plan.id, plan.name, plan.description FROM plan, plan_management WHERE plan_management.IndiOrDepart = "D" AND plan_management.foreign_id = ? AND plan_management.plan_id=plan.id AND plan.current_status = "D"`
)

type IDepartmentDAO interface {
	SelectDepartmentByID(ctx context.Context, db *mssqlx.DBs, id uint64) (result *model.Department, err error)
	SelectAllDepartment(ctx context.Context, db *mssqlx.DBs) (result []*model.DepartmentForPlan, err error)
	SelectDepartmentByView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.Department, err error)
	SelectAllUserWorkingInDepartment(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.UserDetail, err error)
	SelectAllUserUsedToWorkInDepartment(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.UserDetail, err error)
	SelectAllOnGoingPlanInDepartment(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error)
	SelectAllDonePlanInDepartment(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error)
	UpdateDepartment(ctx context.Context, db *mssqlx.DBs, id uint64, name, description string, totalSalary uint64) (err error)
	CreateNewDepartment(ctx context.Context, db *mssqlx.DBs, name, description string, totalSalary uint64) (err error)
}

type departmentDAO struct{}

func (c *departmentDAO) SelectDepartmentByID(ctx context.Context, db *mssqlx.DBs, id uint64) (result *model.Department, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	result = &model.Department{}
	if err = db.GetContext(ctx, result, sqlDepartmentSelectByID, id); err == sql.ErrNoRows {
		result, err = nil, nil
		return
	}

	return
}

func (c *departmentDAO) CreateNewDepartment(ctx context.Context, db *mssqlx.DBs, name, description string, totalSalary uint64) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	if _, er := db.ExecContext(ctx, sqlDepartmentSave, name, description, totalSalary); er != nil {
		err = er
		return
	}
	return
}

func (c *departmentDAO) UpdateDepartment(ctx context.Context, db *mssqlx.DBs, id uint64, name, description string, totalSalary uint64) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	if totalSalary == 0 {
		if _, er := db.ExecContext(ctx, sqlDepartmentUpdateWithoutSalary, name, description, time.Now(), id); er != nil {
			err = er
			return
		}
		return
	}

	if _, er := db.ExecContext(ctx, sqlDepartmentUpdate, name, description, totalSalary, time.Now(), id); er != nil {
		err = er
		return
	}
	return
}

func (c *departmentDAO) SelectAllDepartment(ctx context.Context, db *mssqlx.DBs) (result []*model.DepartmentForPlan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlDepartmentSelectAll)
	return
}

func (c *departmentDAO) SelectDepartmentByView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.Department, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlDepartmentSelectByView, offset, limit)
	return
}

func (c *departmentDAO) SelectAllUserWorkingInDepartment(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.UserDetail, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectWorkingUserInDepartment, id)
	return
}

func (c *departmentDAO) SelectAllUserUsedToWorkInDepartment(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.UserDetail, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectUsedToWorkUserInDepartment, id)
	return
}

func (c *departmentDAO) SelectAllOnGoingPlanInDepartment(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectOnGoingPlanInDepartment, id)
	return
}

func (c *departmentDAO) SelectAllDonePlanInDepartment(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectDonePlanInDepartment, id)
	return
}

var dDao IDepartmentDAO = &departmentDAO{}

// GetUserDAO get user dao
func GetDepartmentDAO() IDepartmentDAO {
	return dDao
}

// SetUserDAO set campaign dao
// NOTE: USE THIS FUNC FOR MOCKING
func SetDepartmentDAO(v IDepartmentDAO) {
	dDao = v
}
