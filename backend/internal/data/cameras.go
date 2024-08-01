package data

import (
	"context"
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
	ModelNo    string    `json:"model_no"`    // String with camera model number/name
	Version    int32     `json:"version"`     // record version
}

type CameraModel struct {
	DB *sql.DB
}

// Insert creates a camera in database
func (c CameraModel) Insert(camera *Camera) error {
	query := `
		insert into cameras (name, mac_address, site_name, model_no)
		values ($1, $2, $3, $4)
		returning id, created_at, version
	`
	args := []any{camera.Name, camera.MacAddress, camera.SiteName, camera.ModelNo}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&camera.ID, &camera.CreatedAt, &camera.Version)
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

	// Create context to terminate long sql queries
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// defer cancel() to ensure context is cancelled prior to Get() returning
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
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

// GetAll() retrieves all cameras from database
func (c CameraModel) GetAll(name string, mac_address string, model_no string, site_name string, filters Filters) ([]*Camera, error) {
	query := `
		select id, created_at, name, mac_address, site_name, model_no, version
		from cameras
		order by id
	`
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cameras := []*Camera{}

	for rows.Next() {
		var camera Camera
		err := rows.Scan(
			&camera.ID,
			&camera.CreatedAt,
			&camera.Name,
			&camera.MacAddress,
			&camera.ModelNo,
			&camera.SiteName,
			&camera.Version,
		)
		if err != nil {
			return nil, err
		}
		cameras = append(cameras, &camera)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return cameras, nil
}

// Update updates a camera in database
func (c CameraModel) Update(camera *Camera) error {
	query := `
		UPDATE cameras
		SET name = $1, mac_address = $2, site_name = $3, model_no = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`
	args := []any{camera.Name, camera.MacAddress, camera.SiteName, camera.ModelNo, camera.ID, camera.Version}

	// Create context to terminate long sql queries
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&camera.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// Delete removes a camera entry from database
func (c CameraModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM cameras
		WHERE id = $1
	`
	// Create context to terminate long sql queries
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	res, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func ValidateCamera(v *validator.Validator, camera *Camera) {
	v.Check(camera.Name != "", "name", "must be provided")
	v.Check(len(camera.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(len(camera.MacAddress) == 12, "mac_address", "must be 12 characters")
	v.Check(validator.Matches(camera.SiteName, siteNameRxp), "site_name", "must be like 'City-Street_Number-Office_Type'")
}
