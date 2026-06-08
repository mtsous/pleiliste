package core

import (
	"context"
)

type Track struct {
	ID      string
	Name    string
	Artists Artists
}

type Tracks []*Track

type Artist struct {
	ID   string
	Name string
}

type Artists []*Artist

type TrackService interface {
	SearchByName(
		ctx context.Context,
		sessionID string,
		name string,
	) (*Tracks, error)
}
