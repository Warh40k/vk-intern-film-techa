package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Type    string `json:"type,omitempty" example:"POST localhost:8080/api/v1/actors/1"`
	Title   string `json:"title,omitempty" example:"input error"`
	Status  int    `json:"status,omitempty" example:"400"`
	Detail  string `json:"detail,omitempty" example:"Failed to get film id. Please, check your input"`
	Message string `json:"-"`
}

func newErrResponse(log *slog.Logger, w http.ResponseWriter, status int, errtype, title, detail, logMessage string) {
	resp := errorResponse{
		Type:    errtype,
		Title:   title,
		Detail:  detail,
		Status:  status,
		Message: logMessage,
	}

	strResp, _ := json.Marshal(resp)

	log.With(slog.String("response", string(strResp)), slog.String("err", logMessage)).Error(resp.Title)

	w.WriteHeader(status)
	w.Write(strResp)
}
