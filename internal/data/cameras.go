package data

import (
	"regexp"
	"time"

	"github.com/chefgoldbloom/pnctool/internal/validator"
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
	Username   string    `json:"username"`    // camera username for admin account
	Password   string    `json:"password"`    // Plaintext password, future: repo integration
	// IpAddress?
}

func ValidateCamera(v *validator.Validator, camera *Camera) {
	v.Check(camera.Name != "", "name", "must be provided")
	v.Check(len(camera.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(len(camera.MacAddress) == 12, "mac_address", "must be 12 characters")
	v.Check(validator.Matches(camera.SiteName, siteNameRxp), "site_name", "must be like 'City-Street_Number-Office_Type'")
}
