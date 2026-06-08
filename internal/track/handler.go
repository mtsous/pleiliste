package track

import (
	"log/slog"
	"net/http"

	"github.com/mtsous/pleiliste/internal/core"
	"github.com/mtsous/pleiliste/internal/util"
)

type handler struct {
	service core.TrackService
}

func NewHandler(service core.TrackService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) SearchByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionID, err := util.GetSessionID(r)
	if err != nil {
		msg := err.Error()
		util.Resp(w, http.StatusInternalServerError, msg)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		slog.Error("failed to search track without required query param name")
		util.Resp(w, http.StatusBadRequest, "query param name is required")
		return
	}

	tracks, err := h.service.SearchByName(ctx, sessionID, name)
	if err != nil {
		msg := err.Error()
		util.Resp(w, http.StatusInternalServerError, msg)
		return
	}

	util.RespRaw(w, http.StatusOK, tracks)
}
