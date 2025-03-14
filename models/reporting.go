package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Report struct {
	bun.BaseModel     `bun:"table:reports"`
	ID                int64     `bun:"id,pk,autoincrement" json:"id"`
	ReportName        string    `bun:"report_name,notnull" json:"report_name"`
	Location          string    `bun:"location,notnull" json:"location"`
	Severity          string    `bun:"severity,notnull" json:"severity"`
	ReportedBy        string    `bun:"reported_by,notnull" json:"reported_by"`
	Status            string    `bun:"status,notnull,default:'Pending'" json:"status"`
	DateReported      time.Time `bun:"date_reported,notnull,default:current_timestamp" json:"date_reported"`
	ActionDescription string    `bun:"action_description" json:"action_description"`
	PhotoUrls         []string  `bun:"photo_urls,array" json:"photo_urls"`
	Department        string    `bun:"department,notnull" json:"department"`
}

type FireReport struct {
	Report
}

type EMSReport struct {
	Report
}

type AVSReport struct {
	Report
}
