package model

import (
	"strconv"
	"time"

	libs "DBproject1/lib"
)

// User model user
type User struct {
	ID         uint64
	Username   string
	PwdCounter int `db:"pwd_counter"`
	Status     byte
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Checksum   uint64
}

type UserDetail struct {
	ID           uint64    `json:"ID"`
	Name         string    `json:"Name"`
	DOB          string    `db:"DOB" json:"DOB"`
	Sex          string    `json:"Sex"`
	Position     string    `json:"Position"`
	PhoneNum     string    `json:"PhoneNum"`
	NationalID   string    `db:"national_id" json:"NationalID"`
	Salary       int       `json:"Salary"`
	Username     string    `json:"Username"`
	DepartmentID uint64    `json:"DepartmentID"`
	Password     string    `json:"Password"`
	CreatedAt    time.Time `db:"created_at" json:"CreatedAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"UpdatedAt"`
}

type UserDetailString struct {
	ID         uint64 `json:"ID"`
	Name       string `json:"Name"`
	DOB        string `db:"DOB" json:"DOB"`
	Sex        string `json:"Sex"`
	Position   string `json:"Position"`
	PhoneNum   string `json:"PhoneNum"`
	NationalID string `db:"national_id" json:"NationalID"`
	Salary     string `json:"Salary"`
}

type Department struct {
	ID          uint64    `json:"ID"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	TotalSalary uint64    `db:"total_salary" json:"TotalSalary"`
	CreatedAt   time.Time `db:"created_at" json:"CreatedAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"UpdatedAt"`
}

type Plan struct {
	ID            uint64    `json:"ID"`
	Name          string    `json:"Name"`
	Description   string    `json:"Description"`
	FatherPlanID  int       `json:"FatherPlanID" db:"father_plan_id"`
	CurrentStatus string    `json:"CurrentStatus" db:"current_status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type PlanForUI struct {
	ID           uint64 `json:"ID"`
	Name         string `json:"Name"`
	Description  string `json:"Description"`
	FatherPlanID string `json:"FatherPlanID"`
}

type Plant_Management struct {
	IndiOrDepart string
	ForeignID    uint64
	PlanID       uint64
}

type UserDetailForPlan struct {
	ID       uint64 `json:"ID"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
}

type DepartmentForPlan struct {
	ID   uint64 `json:"ID"`
	Name string `json:"Name"`
}

type DepartmentManagement struct {
	DepartmentID uint64    `db:"department_id"`
	UserID       uint64    `db:"user_id"`
	Status       int       `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UserPlan struct {
	DonePlan    []*Plan `json:"DonePlan"`
	OnGoingPlan []*Plan `json:"OnGoingPlan"`
}

type UserDepartment struct {
	OldUser []*UserDetail `json:"OldUser"`
	NewUser []*UserDetail `json:"NewUser"`
}

// UserWithRP User with RP for UI
type UserWithRP struct {
	Username    string     `json: username`
	Groups      []string   `json: groups`
	Permissions [][]string `json: permissions`
}

// GroupWithPermissions group of user
type GroupWithPermissions struct {
	GroupName   string     `json:name`
	Permissions [][]string `json: permissions`
}

// Sum calculate sip hash sum
func (c User) Sum(k0, k1 uint64) uint64 {
	sum := libs.ConcatCopyPreAllocate([][]byte{
		[]byte(strconv.FormatUint(c.ID, 10)),
		{c.Status},
		[]byte(c.Username),
	})

	return uint64(libs.SipHash48(k0, k1, []byte(sum)))
}

// ValidateChecksum validate record checksum
func (c User) ValidateChecksum(k0, k1 uint64) bool {
	return c.Sum(k0, k1) == c.Checksum
}

// UserSecurity user security table
type UserSecurity struct {
	Username  string
	Gr        uint64
	Role      uint64
	Password  []byte
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Checksum  uint64
}

// Sum do sum on user security
func (c *UserSecurity) Sum(k0, k1 uint64) uint64 {
	sum := libs.ConcatCopyPreAllocate([][]byte{
		[]byte(c.Username),
		libs.Uint64ToBytes(c.Gr),
		libs.Uint64ToBytes(c.Role),
		c.Password,
	})

	return uint64(libs.SipHash48(k0, k1, []byte(sum)))
}

// ValidateChecksum validate checksum
func (c *UserSecurity) ValidateChecksum(k0, k1 uint64) bool {
	return c.Sum(k0, k1) == c.Checksum
}
