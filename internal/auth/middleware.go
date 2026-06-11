package auth

import (
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/mtsous/pleiliste/internal/session"
)

func SessionMiddleware(store *session.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("hello world")

			switch r.URL.Path {
			case
				"/auth/spotify",
				"/auth/spotify/callback":
				next.ServeHTTP(w, r)
			}

			cookie, err := r.Cookie("session_id")
			if err == http.ErrNoCookie {
				slog.Info("hello world 1")

				redirectSpotifyAuth(w, r)
				return
			}

			sess, err := store.Get(cookie.Value)
			if err != nil {
				slog.Info("hello world 2")

				redirectSpotifyAuth(w, r)
				return
			}

			if sess.ExpiresAt.Before(time.Now()) {
				slog.Info("hello world 3")

				redirectSpotifyAuth(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})

	}
}

func redirectSpotifyAuth(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	q := url.Values{"state": {state}}
	slog.Info("mahmoud almadinejad")
	http.Redirect(w, r, "/auth/spotify?"+q.Encode(), http.StatusTemporaryRedirect)
}
