package models

import "time"

// Proposal
type Proposal struct {
	// Collection of changes
	ID          int
	Title       string
	Description string
	CreatedAt   time.Time
	Changes     []Change
	Submitted   bool `default:"false"`
}

type Change interface{}
