package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Email      string   `json:"email"`
	Roles      []string `json:"roles"`
	RoleIDs    []string `json:"rolesId"`
	RolesLower []string `json:"rolesLower"`
	Name       string   `json:"name"`
	Username   string   `json:"username"`
	jwt.RegisteredClaims
}

type Date struct {
	time.Time
}

type DateRange struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format("2006-01-02"))
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	t, ok := value.(time.Time)
	if !ok {
		return errors.New("value is not a time.Time")
	}

	d.Time = t
	return nil
}

func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}

	return d.Time, nil
}
