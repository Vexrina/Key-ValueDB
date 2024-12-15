package utils

import (
	"BD/pkg/xlog"
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteApiOK(w http.ResponseWriter, answer any, status int) {
	if answer != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(answer); err != nil {
			xlog.Error("writeApiOK error", xlog.ErrorField(err))
			WriteApiError(w, fmt.Sprintf("got error during writing answer, %s", err.Error()), http.StatusInternalServerError)
		}
	}

	w.WriteHeader(status)
}

func WriteApiError(w http.ResponseWriter, msg string, status int) {
	xlog.Error("some method got error", xlog.Field("error", msg))
	http.Error(w, msg, status)
}
