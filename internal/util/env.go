package util

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/mtsous/pleiliste/internal/core"
)

func GetEnv() *core.Env {
	godotenv.Load()

	return &core.Env{
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}
}
