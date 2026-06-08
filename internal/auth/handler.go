package auth

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mtsous/pleiliste/internal/core"
	"github.com/mtsous/pleiliste/internal/util"
)

type handler struct {
	spotify      core.SpotifyClient
	sessionStore core.SessionStore
}

func NewHandler(
	spotify core.SpotifyClient,
	sessionStore core.SessionStore,
) *handler {
	return &handler{
		spotify:      spotify,
		sessionStore: sessionStore,
	}
}

func (h *handler) HandleSpotify(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		slog.Error("state query param not provided for spotify auth")
		util.Resp(w, http.StatusBadRequest, "state is required")
		return
	}

	url := h.spotify.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *handler) HandleSpotifyCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	tokensResp, err := h.spotify.GetTokens(ctx, code)
	if err != nil {
		util.Resp(w, http.StatusInternalServerError, "failed to authenticate")
		return
	}

	sessionID := uuid.New().String()
	expiresInSeconds := time.Second * time.Duration(tokensResp.ExpiresIn)
	expiresAt := time.Now().Add(expiresInSeconds)

	h.sessionStore.Set(sessionID, &core.Session{
		State:        state,
		AccessToken:  tokensResp.AccessToken,
		RefreshToken: tokensResp.RefreshToken,
		ExpiresAt:    expiresAt,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 7,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
