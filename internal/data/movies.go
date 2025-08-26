package data

import (
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`

	Title   string   `json:"title"`
	Year    int32    `json:"year,omitzero"`
	Runtime Runtime  `json:"runtime,omitzero"`
	Genres  []string `json:"genres,omitzero"`

	// The version number starts at 1 and will be incremented
	// each time the movie information is updated
	Version int32 `json:"version"`
}
