package auth

import (
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/mtsous/pleiliste/internal/session"
)

func SessionMiddleware(store *session.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err == http.ErrNoCookie {
				redirectSpotifyAuth(w, r)
				return
			}

			if _, err := store.Get(cookie.Value); err != nil {
				redirectSpotifyAuth(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})

	}
}

func redirectSpotifyAuth(w http.ResponseWriter, r *http.Request) {
	code := uuid.New().String()
	q := url.Values{"code": {code}}
	http.Redirect(w, r, "/auth/spotify"+q.Encode(), http.StatusTemporaryRedirect)
}
