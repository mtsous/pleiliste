package core

import (
	"context"
	"errors"
)

type SpotifyTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type SpotifyTracks struct {
	Tracks struct {
		Items []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"items"`
	} `json:"tracks"`
}

type SpotifyClient interface {
	GetAuthURL(state string) string
	GetTokens(
		ctx context.Context,
		code string,
	) (*SpotifyTokens, error)
	SearchByName(
		ctx context.Context,
		sessionID string,
		name string,
	) (*Tracks, error)
}

var (
	ErrSpotifyStatusCodeNotOK = errors.New("spotify response status was 200 ok")
	ErrSpotifySearchError     = errors.New("spotify search error")
	ErrSpotifyNoTracksFound   = errors.New("spotify did not found any tracks")
)
