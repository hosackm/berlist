package event

import (
	"time"
)

type Event struct {
	Date     time.Time `json:"date"`
	Name     string    `json:"name"`
	Venue    string    `json:"venue"`
	Link     string    `json:"link"`
	Subtitle string    `json:"subtitle,omitempty"`
}
