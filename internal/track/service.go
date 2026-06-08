package track

import (
	"context"

	"github.com/mtsous/pleiliste/internal/core"
)

type service struct {
	spotify core.SpotifyClient
}

func NewService(spotify core.SpotifyClient) *service {
	return &service{spotify: spotify}
}

func (s *service) SearchByName(
	ctx context.Context,
	sessionID string,
	name string,
) (*core.Tracks, error) {
	return s.spotify.SearchByName(ctx, sessionID, name)
}
