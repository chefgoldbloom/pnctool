package data

import (
	"database/sql"
	"errors"
	"regexp"
	"time"

	"github.com/chefgoldbloom/pnctool/backend/internal/validator"
)

var (
	siteNameRxp = regexp.MustCompile(".*-.*-(OPS|COE|GLH)$")
)

type Camera struct {
	ID         int64     `json:"id"`          // Unique integer ID for the camera
	CreatedAt  time.Time `json:"created_at"`  // Timestamp for when the camera is added to our database
	Name       string    `json:"name"`        // Camera Name
	MacAddress string    `json:"mac_address"` // Serial number for camera like 'ACC...'
	SiteName   string    `json:"site_name"`   // String with site name, future: foreign key linked to Site table
	Username   string    `json:"-"`           // camera username for admin account
	Password   string    `json:"-"`           // Plaintext password, future: repo integration
	ModelNo    string    `json:"model_no"`
	// IpAddress?
}

type CameraModel struct {
	DB *sql.DB
}

// Insert creates a camera in database
func (c CameraModel) Insert(camera *Camera) error {
	query := `
		insert into cameras (name, mac_address, site_name, model_no)
		values ($1, $2, $3, $4)
		returning id, created_at
	`
	args := []any{camera.Name, camera.MacAddress, camera.SiteName, camera.ModelNo}

	return c.DB.QueryRow(query, args...).Scan(&camera.ID, &camera.CreatedAt)
}

// Get retrieves camera from database
func (c CameraModel) Get(id int64) (*Camera, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		select id, created_at, name, mac_address, site_name, model_no
		from cameras
		where id = $1
	`
	var camera Camera

	err := c.DB.QueryRow(query, id).Scan(
		&camera.ID,
		&camera.CreatedAt,
		&camera.Name,
		&camera.MacAddress,
		&camera.SiteName,
		&camera.ModelNo,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &camera, nil
}

// Update updates a camera in database
func (c CameraModel) Update(camera *Camera) error {
	return nil
}

// Delete removes a camera entry from database
func (c CameraModel) Delete(id int64) error {
	return nil
}

func ValidateCamera(v *validator.Validator, camera *Camera) {
	v.Check(camera.Name != "", "name", "must be provided")
	v.Check(len(camera.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(len(camera.MacAddress) == 12, "mac_address", "must be 12 characters")
	v.Check(validator.Matches(camera.SiteName, siteNameRxp), "site_name", "must be like 'City-Street_Number-Office_Type'")
}
