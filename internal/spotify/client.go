package spotify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/mtsous/pleiliste/internal/core"
)

var (
	redirectURI        = "http://127.0.0.1:8080/auth/spotify/callback"
	spotifyAPIURL      = "https://api.spotify.com/v1"
	spotifyAccountsURL = "https://accounts.spotify.com"
)

type Client struct {
	env        *core.Env
	httpClient *http.Client
	store      core.SessionStore
}

func NewClient(
	env *core.Env,
	httpClient *http.Client,
	sessionStore core.SessionStore,
) *Client {
	return &Client{
		env:        env,
		httpClient: httpClient,
		store:      sessionStore,
	}
}

func (c *Client) GetAuthURL(state string) string {
	u, _ := url.Parse(spotifyAccountsURL + "/authorize")
	u.RawQuery = url.Values{
		"client_id":     {c.env.SpotifyClientID},
		"response_type": {"code"},
		"redirect_uri":  {redirectURI},
		"state":         {state},
		"scope":         {"playlist-modify-public"},
		"show_dialog":   {"true"},
	}.Encode()

	return u.String()
}

func (c *Client) GetTokens(
	ctx context.Context,
	code string,
) (*core.SpotifyTokens, error) {
	form := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {redirectURI},
	}.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		spotifyAccountsURL+"/api/token",
		strings.NewReader(form),
	)
	if err != nil {
		slog.Error(
			"error creating spotify tokens request",
			"err", err.Error(),
		)
		return nil, err
	}

	basicTokenString := []byte(c.env.SpotifyClientID + ":" + c.env.SpotifyClientSecret)
	basicTokenBase64 := base64.StdEncoding.EncodeToString(basicTokenString)

	req.Header.Set("Authorization", "Basic "+basicTokenBase64)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error(
			"error requesting spotify tokens",
			"err", err.Error(),
		)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error(
			"error status code on spotify tokens response",
			"body", string(body),
			"status", resp.StatusCode,
		)
		return nil, core.ErrSpotifyStatusCodeNotOK
	}

	var tokensResp *core.SpotifyTokens

	err = json.NewDecoder(resp.Body).Decode(&tokensResp)
	if err != nil {
		slog.Error(
			"error decoding spotify tokens",
			"err", err.Error(),
		)
		return nil, err
	}

	return tokensResp, nil
}

func (c *Client) SearchByName(
	ctx context.Context,
	sessionID string,
	name string,
) (*core.Tracks, error) {
	session, err := c.store.Get(sessionID)
	if err != nil {
		slog.Error(
			"failed to get session store",
			"session_id", sessionID,
		)
		return nil, err
	}

	params := url.Values{
		"q":    {name},
		"type": {"track"},
	}.Encode()

	url := spotifyAPIURL + "/search?" + params

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		slog.Error(
			"failed to create spotify search request",
			"query", name,
			"err", err,
		)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+session.AccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error(
			"failed to request spotify search",
			"query", name,
			"err", err,
		)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error(
			"unexpected response status requesting spotify search",
			"body", string(body),
			"err", err,
		)
		return nil, core.ErrSpotifySearchError
	}

	return c.decodeTracks(resp.Body)
}

func (c *Client) decodeTracks(body io.ReadCloser) (*core.Tracks, error) {
	var response core.SpotifyTracks

	err := json.NewDecoder(body).Decode(&response)
	if err != nil {
		slog.Error("could not decode body for track search response")
		return nil, err
	}

	numberOfTracks := len(response.Tracks.Items)

	if numberOfTracks <= 0 {
		slog.Error("no tracks found for search query")
		return nil, core.ErrSpotifyNoTracksFound
	}

	result := make(core.Tracks, numberOfTracks)
	for i, v := range response.Tracks.Items {
		artists := make(core.Artists, len(v.Artists))

		track := core.Track{
			ID:   v.ID,
			Name: v.Name,
			ISRC: v.ExternalIDs.ISRC,
		}

		for j, vv := range v.Artists {
			artist := core.Artist{
				ID:   vv.ID,
				Name: vv.Name,
			}
			artists[j] = &artist
		}

		track.Artists = artists

		result[i] = &track
	}

	return &result, nil
}
