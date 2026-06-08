package util

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type messageResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func Resp(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(messageResponse{
		Status: status,
		Msg:    msg,
	})
}

func RespRaw(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(v)
}

func GetSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		slog.Error("failed to get session_id cookie")
		return "", err
	}

	return cookie.Value, nil
}
