package attendance

import "time"

// Model is the base for all models.
type Model struct {
	// ID uniquely identifies an entity.
	ID uint64
	// CreatedAt is the date the entity was created at.
	CreatedAt time.Time
	// UpdatedAt is the date the entity was updated at.
	UpdatedAt time.Time
}

