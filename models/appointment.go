package models

import (
	"time"

	"github.com/uptrace/bun"
)

type PriorityEnum string

const (
	PriorityLow    PriorityEnum = "low"
	PriorityMedium PriorityEnum = "medium"
	PriorityHigh   PriorityEnum = "high"
)

type Appointment struct {
	bun.BaseModel `bun:"table:appointments"`

	ID              int64          `bun:"id,pk,autoincrement" json:"id"`
	VisitorName     string         `bun:"visitor_name,notnull" json:"visitor_name"`
	Purpose         string         `bun:"purpose,notnull" json:"purpose"`
	WhoToSee        string         `bun:"who_to_see,notnull" json:"who_to_see"`
	Department      DepartmentEnum `bun:"department,notnull" json:"department"`
	AppointmentDate time.Time      `bun:"appointment_date,notnull" json:"appointment_date"`
	TimeIn          time.Time      `bun:"time_in,notnull" json:"time_in"`
	TimeOut         time.Time      `bun:"time_out,notnull" json:"time_out"`
	Priority        PriorityEnum   `bun:"priority,notnull" json:"priority"`
	Notes           string         `bun:"notes" json:"notes"`

	CreatedAt time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}
