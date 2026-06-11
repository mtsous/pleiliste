package track

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/mtsous/pleiliste/internal/core"
	"github.com/mtsous/pleiliste/internal/util"
)

type handler struct {
	tmpl    *template.Template
	service core.TrackService
}

func NewHandler(
	tmpl *template.Template,
	service core.TrackService,
) *handler {
	return &handler{
		tmpl:    tmpl,
		service: service,
	}
}

type IndexPage struct {
	Title string
}

func (h *handler) SearchIndex(w http.ResponseWriter, r *http.Request) {
	data := IndexPage{
		Title: "Escolha uma música",
	}

	h.tmpl.ExecuteTemplate(w, "index.html", data)
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
		slog.Error("error query param name is required")
		util.Resp(w, http.StatusBadRequest, "query param name is required")
		return
	}

	tracks, err := h.service.SearchByName(ctx, sessionID, name)
	if err != nil {
		msg := err.Error()
		util.Resp(w, http.StatusInternalServerError, msg)
		return
	}

	h.tmpl.ExecuteTemplate(w, "tracks", tracks)
}
