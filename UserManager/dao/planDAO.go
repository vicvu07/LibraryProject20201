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
	sqlSelectAllPlan             = "SELECT * FROM plan"
	sqlSelectByViewPlan          = "SELECT * FROM plan ORDER BY id DESC LIMIT ?, ?"
	sqlSelectByIDPlan            = "SELECT * FROM plan WHERE id = ?"
	sqlSelectAllFinishedPlan     = `SELECT * FROM plan WHERE current_status = "D"`
	sqlSelectAllOnGoingPlan      = `SELECT * FROM plan WHERE current_status = "O"`
	sqlSelectChildrenOnGoingPlan = `SELECT * FROM plan WHERE father_plan_id = ? AND current_status = "O"`
	sqlSelectChildrenDonePlan    = `SELECT * FROM plan WHERE father_plan_id = ? AND current_status = "D"`
	sqlSelectUserInPlan          = `SELECT user_detail.id, user_detail.name, user_detail.username FROM user_detail, plan_management WHERE plan_management.IndiOrDepart = "I" and plan_management.plan_id = ? and user_detail.id = plan_management.foreign_id`
	sqlSelectDepartmentInPlan    = `SELECT department.id, department.name FROM department, plan_management WHERE plan_management.IndiOrDepart = "D" and plan_management.plan_id = ? and department.id = plan_management.foreign_id`
	sqlDeleteOneFromPlan         = "DELETE FROM plan_management WHERE IndiOrDepart = ?, plan_id = ? AND foreign_id = ?"
	sqlDeleteAllFromPlan         = "DELETE FROM plan_management WHERE plan_id = ? AND IndiOrDepart = ?"
	sqlSavePlan                  = `INSERT INTO plan(name,description,father_plan_id,current_status) VALUE (?,?,?,"O")`
	sqlUpdatePlan                = "UPDATE plan SET name = ?, description = ?, father_plan_id = ? WHERE id = ?"
	sqlDonePlan                  = `UPDATE plan SET current_status = "D", updated_at = ? WHERE id = ?`
	sqlAddToPlan                 = "INSERT INTO plan_management(IndiOrDepart,foreign_id,plan_id) VALUES (?,?,?)"
)

type IPlanDAO interface {
	SelectPlanByView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.Plan, err error)
	SelectPlanByID(ctx context.Context, db *mssqlx.DBs, id uint64) (result *model.Plan, err error)
	SelectAllFinishedPlan(ctx context.Context, db *mssqlx.DBs) (result []*model.Plan, err error)
	SelectAllOnGoingPlan(ctx context.Context, db *mssqlx.DBs) (result []*model.Plan, err error)
	SelectChildrenOnGoingPlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error)
	SelectChildrenDonePlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error)
	SelectUserInPlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.UserDetailForPlan, err error)
	SelectDepartmentInPlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.DepartmentForPlan, err error)
	SavePlan(ctx context.Context, db *mssqlx.DBs, name, description string, fatherPlanID uint64) (err error)

	DeleteOneFromPlan(ctx context.Context, db *mssqlx.DBs, InviOrDepart string, id uint64, foreignID uint64) (err error)
	DeleteAllFromPlan(ctx context.Context, db *mssqlx.DBs, id uint64, InviOrDepart string) (err error)

	UpdatePlan(ctx context.Context, db *mssqlx.DBs, id uint64, name, description string, fatherPlanID uint64) (err error)
	DonePlan(ctx context.Context, db *mssqlx.DBs, id uint64) (err error)
	AddToPlan(ctx context.Context, db *mssqlx.DBs, IndiOrDepart string, foreignID, planID uint64) (err error)
}

type planDAO struct{}

func (c *planDAO) SelectPlanByView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.Plan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectByViewPlan, offset, limit)
	return
}

func (c *planDAO) SelectAllFinishedPlan(ctx context.Context, db *mssqlx.DBs) (result []*model.Plan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectAllFinishedPlan)
	return
}

func (c *planDAO) SelectAllOnGoingPlan(ctx context.Context, db *mssqlx.DBs) (result []*model.Plan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectAllOnGoingPlan)
	return
}
func (c *planDAO) SelectPlanByID(ctx context.Context, db *mssqlx.DBs, id uint64) (result *model.Plan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	result = &model.Plan{}

	if err = db.GetContext(ctx, result, sqlSelectByIDPlan, id); err == sql.ErrNoRows {
		result, err = nil, nil
		return
	}

	return
}

func (c *planDAO) SavePlan(ctx context.Context, db *mssqlx.DBs, name, description string, fatherPlanID uint64) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	if _, er := db.ExecContext(ctx, sqlSavePlan, name, description, fatherPlanID); er != nil {
		err = er
		return
	}
	return
}

func (c *planDAO) UpdatePlan(ctx context.Context, db *mssqlx.DBs, id uint64, name, description string, fatherPlanID uint64) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	if _, er := db.ExecContext(ctx, sqlUpdatePlan, name, description, fatherPlanID, id); er != nil {
		err = er
		return
	}
	return
}

func (c *planDAO) DonePlan(ctx context.Context, db *mssqlx.DBs, id uint64) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	if _, er := db.ExecContext(ctx, sqlDonePlan, time.Now(), id); er != nil {
		err = er
		return
	}
	return
}

func (c *planDAO) AddToPlan(ctx context.Context, db *mssqlx.DBs, IndiOrDepart string, foreignID, planID uint64) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	if _, er := db.ExecContext(ctx, sqlAddToPlan, IndiOrDepart, foreignID, planID); er != nil {
		err = er
		return
	}
	return
}

func (c *planDAO) SelectChildrenOnGoingPlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectChildrenOnGoingPlan, id)
	return
}

func (c *planDAO) SelectChildrenDonePlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.Plan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectChildrenDonePlan, id)
	return
}

func (c *planDAO) SelectUserInPlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.UserDetailForPlan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectUserInPlan, id)
	return
}

func (c *planDAO) SelectDepartmentInPlan(ctx context.Context, db *mssqlx.DBs, id uint64) (result []*model.DepartmentForPlan, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlSelectDepartmentInPlan, id)
	return
}

func (c *planDAO) DeleteOneFromPlan(ctx context.Context, db *mssqlx.DBs, InviOrDepart string, id uint64, foreignID uint64) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	if _, er := db.ExecContext(ctx, sqlDeleteOneFromPlan, InviOrDepart, id, foreignID); er != nil {
		err = er
		return
	}
	return
}

func (c *planDAO) DeleteAllFromPlan(ctx context.Context, db *mssqlx.DBs, id uint64, InviOrDepart string) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	if _, er := db.ExecContext(ctx, sqlDeleteAllFromPlan, id, InviOrDepart); er != nil {
		err = er
		return
	}
	return
}

var pDAO IPlanDAO = &planDAO{}

// GetUserDAO get user dao
func GetPlanDAO() IPlanDAO {
	return pDAO
}

// SetUserDAO set campaign dao
// NOTE: USE THIS FUNC FOR MOCKING
func SetPlanDAO(v IPlanDAO) {
	pDAO = v
}
