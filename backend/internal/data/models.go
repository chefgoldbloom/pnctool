package data

import (
	"database/sql"
	"errors"
)

// Define ErrRecordNotFound.  Will return this from Get() method when
// looking up a camera that doesn't exist
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Create a Models struct that wraps CameraModel
type Models struct {
	Cameras CameraModel
}

// Create a New() method that will instantiate Models
func NewModels(db *sql.DB) Models {
	return Models{
		Cameras: CameraModel{DB: db},
	}
}
