package models

import (
	"time"

	"github.com/uptrace/bun"
)

type SeverityEnum string

const (
	SeverityLow      SeverityEnum = "Low"
	SeverityModerate SeverityEnum = "Moderate"
	SeverityHigh     SeverityEnum = "High"
	SeverityCritical SeverityEnum = "Critical"
)

type IncidentTypeEnum string

const (
	IncidentFire     IncidentTypeEnum = "Fire"
	IncidentMedical  IncidentTypeEnum = "Medical"
	IncidentSecurity IncidentTypeEnum = "Security"
	IncidentOther    IncidentTypeEnum = "Other"
)

type Incident struct {
	bun.BaseModel `bun:"table:incidents"`

	ID                int64            `bun:"id,pk,autoincrement" json:"id"`
	AgentID           string           `bun:"agent_id,notnull" json:"agent_id"`
	Department        DepartmentEnum   `bun:"department,notnull" json:"department"`
	IncidentType      IncidentTypeEnum `bun:"incident_type,notnull" json:"incident_type"`
	Severity          SeverityEnum     `bun:"severity,notnull" json:"severity"`
	CallerFullName    string           `bun:"caller_full_name,notnull" json:"caller_full_name"`
	CallerPhoneNumber string           `bun:"caller_phone_number,notnull" json:"caller_phone_number"`
	CallerLocation    string           `bun:"caller_location,notnull" json:"caller_location"`
	PeopleInvolved    int              `bun:"people_involved,notnull" json:"people_involved"`
	IncidentReport    string           `bun:"incident_report,notnull" json:"incident_report"`
	StaffID           int64            `bun:"staff_id,notnull" json:"staff_id"`

	Staff *Staff `bun:"rel:belongs-to,join:staff_id=id" json:"staff"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}
