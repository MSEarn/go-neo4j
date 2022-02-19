package http_helper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func Decode(r *http.Request, data interface{}) error {
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(data)
}

func RespondOK(w http.ResponseWriter, code uint64, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := &ResponseData{
		Response: Response{
			Code:    code,
			Message: msg,
		},
		Data: data,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Unable to write response, err: %v", err))
	}
}
