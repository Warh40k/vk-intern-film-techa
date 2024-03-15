package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Status  int    `json:"status,omitempty"`
	Detail  string `json:"detail,omitempty"`
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

	log.With(slog.String("resp", fmt.Sprintf("%+v", resp))).Error(resp.Title)

	w.WriteHeader(status)
	w.Write(strResp)
}
