package models

import (
	"time"

	"github.com/uptrace/bun"
)

type PositionEnum string

const (
	PositionSSA        PositionEnum = "SSA"
	PositionDirector   PositionEnum = "Director"
	PositionIT         PositionEnum = "IT"
	PositionCallCenter PositionEnum = "Call Center"
	PositionStaff      PositionEnum = "Staff"
	PositionHR         PositionEnum = "HR"
)

type DepartmentEnum string

const (
	DeptHomelandSecurity DepartmentEnum = "Homeland Security"
	DeptAVS              DepartmentEnum = "AVS"
	DeptEMS              DepartmentEnum = "EMS"
	DeptFireService      DepartmentEnum = "Fire Service"
)

type RoleEnum string

const (
	RoleAdmin    RoleEnum = "Admin"
	RoleSSA      RoleEnum = "SSA"
	RoleDirector RoleEnum = "Director"
	RoleStaff    RoleEnum = "Staff"
)

type Staff struct {
	bun.BaseModel `bun:"table:staff"`

	ID                 int64          `bun:"id,pk,autoincrement" json:"id"`
	FirstName          string         `bun:"first_name,notnull" json:"first_name"`
	MiddleName         string         `bun:"middle_name" json:"middle_name"`
	LastName           string         `bun:"last_name,notnull" json:"last_name"`
	Email              string         `bun:"email,unique,notnull" json:"email"`
	Password           string         `bun:"password,notnull" json:"-"`
	AgentID            string         `bun:"agent_id,unique,notnull" json:"agent_id"`
	ProfilePhoto       string         `bun:"profile_photo" json:"profile_photo"`
	Position           PositionEnum   `bun:"position,notnull" json:"position"`
	Address            string         `bun:"address" json:"address"`
	Department         DepartmentEnum `bun:"department,notnull" json:"department"`
	DateOfBirth        time.Time      `bun:"date_of_birth,notnull" json:"date_of_birth"`
	StateOfOrigin      string         `bun:"state_of_origin,notnull" json:"state_of_origin"`
	Role               RoleEnum       `bun:"role,notnull" json:"role"`
	MustChangePassword bool           `bun:"must_change_password,notnull,default:true" json:"must_change_password"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}
