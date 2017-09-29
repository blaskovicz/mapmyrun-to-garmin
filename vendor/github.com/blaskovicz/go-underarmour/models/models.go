package models

import "time"

type Linkable struct {
	// eg: links.self[0].href
	Links map[string][]map[string]string `json:"_links,omitempty"`
}
type ErrorResponse struct {
	// eg: diagnostics.validation_failures[0].__all__[0]
	Diagnostics map[string][]map[string][]string `json:"_diagnostics,omitempty"`
	Linkable
}

type User struct {
	ID                int       `json:"id"`
	Gender            string    `json:"gender"`
	PreferredLanguage string    `json:"preferred_language"`
	Introduction      string    `json:"introduction"`
	Username          string    `json:"username"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	DisplayName       string    `json:"display_name"`
	LastInitial       string    `json:"last_initial"`
	DateJoined        time.Time `json:"date_joined"`
	ProfileStatement  string    `json:"profile_statement"`
	Hobbies           string    `json:"hobbies"`
	TimeZone          string    `json:"time_zone"`
	GoalStatement     string    `json:"goal_statement"`
	Location          struct {
		Country  string `json:"country"`
		Region   string `json:"region"`
		Locality string `json:"locality"`
	} `json:"location"`
	Linkable
}

type Route struct {
	TotalDescent      float64   `json:"total_descent"`
	TotalAscent       float64   `json:"total_ascent"`
	City              string    `json:"city"`
	DataSource        string    `json:"data_source"`
	Description       string    `json:"description"`
	UpdatedAt         time.Time `json:"updated_datetime"`
	CreatedAt         time.Time `json:"created_datetime"`
	Country           string    `json:"country"`
	StartingPointType string    `json:"starting_point_type"`
	StartingLocation  struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"starting_location"`
	Distance     float64 `json:"distance"`
	Name         string  `json:"name"`
	State        string  `json:"state"`
	MaxElevation float64 `json:"max_elevation"`
	MinElevation float64 `json:"min_elevation"`
	PostalCode   string  `json:"postal_code"`
	Linkable
}
