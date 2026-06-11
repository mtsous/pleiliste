package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mtsous/pleiliste/internal/auth"
	"github.com/mtsous/pleiliste/internal/session"
	"github.com/mtsous/pleiliste/internal/spotify"
	"github.com/mtsous/pleiliste/internal/track"
	"github.com/mtsous/pleiliste/internal/util"
	"github.com/mtsous/pleiliste/web"
)

func main() {
	env := util.GetEnv()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	sessionStore := session.NewStore(3600)

	spotifyClient := spotify.NewClient(env, client, sessionStore)

	trackService := track.NewService(spotifyClient)

	tmpl := web.HandleTmpl()

	authHandler := auth.NewHandler(spotifyClient, sessionStore)
	trackHandler := track.NewHandler(tmpl, trackService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /auth/spotify", authHandler.HandleSpotify)
	mux.HandleFunc("GET /auth/spotify/callback", authHandler.HandleSpotifyCallback)
	mux.HandleFunc("GET /", trackHandler.SearchIndex)
	mux.HandleFunc("GET /tracks", trackHandler.SearchByName)

	handler := auth.SessionMiddleware(sessionStore)(mux)

	web.HandleStatic(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "127.0.0.1:" + port
	log.Printf("server listening on http://%s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
