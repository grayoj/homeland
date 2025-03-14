package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Document struct {
	bun.BaseModel `bun:"table:documents"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Name          string    `bun:"name,notnull"`
	URL           string    `bun:"url,notnull"`
	UploadedBy    string    `bun:"uploaded_by,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:current_timestamp"`
}
