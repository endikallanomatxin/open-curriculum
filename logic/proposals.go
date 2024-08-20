package logic

import "time"

// Proposal
type Proposal struct {
	// Collection of changes
	ID          int64
	Title       string
	Description string
	CreatedAt   time.Time
	Changes     []Change
	Submitted   bool `default:"false"`
}

type Change interface{}
