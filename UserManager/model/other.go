package model

import "time"

type Resource struct {
	ID          uint64    `json:"id" db:"id"`
	Name        string    `json:"Name" db:"name"`
	Description string    `json:"Description" db:"description"`
	CreatedAt   time.Time `json:"CreatedAt" db:"created_at"`
	UpdatedAt   time.Time `json:"UpdatedAt" db:"updated_at"`
}

type Action struct {
	ID          uint64    `json:"id" db:"id"`
	Name        string    `json:"Name" db:"name"`
	Description string    `json:"Description" db:"description"`
	CreatedAt   time.Time `json:"CreatedAt" db:"created_at"`
	UpdatedAt   time.Time `json:"UpdatedAt" db:"updated_at"`
}

type Role struct {
	ID          uint64    `json:"id" db:"id"`
	Name        string    `json:"Name" db:"name"`
	Description string    `json:"Description" db:"description"`
	CreatedAt   time.Time `json:"CreatedAt" db:"created_at"`
	UpdatedAt   time.Time `json:"UpdatedAt" db:"updated_at"`
}

type PermissionForUI struct {
	TPerm [][]string `json:"TPerm"`
}

type RoleForUI struct {
	TRole []string `json:"TRole"`
}
