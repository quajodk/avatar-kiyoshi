package media

import (
	"database/sql"
	"path/filepath"
	"time"
)

// Media
type Media struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Source    string         `json:"source"`
	Thumbnail sql.NullString `json:"thumbnail"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"update_at"`
}

var mediaDirectory string = filepath.Join("uploads")
