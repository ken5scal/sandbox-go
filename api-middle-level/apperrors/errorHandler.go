package apperrors

import (
	"encoding/json"
	"errors"
	mdw "github.com/ken5scal/api-go-mid-level/api/middlewares"
	"log"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *MyAppError
	if !errors.As(err, &appErr) {
		appErr = &MyAppError{ErrCode: Unknown, Message: "internal process failed", Err: err}
	}

	traceID := mdw.GetTraceID(r.Context())
	log.Printf("[%d]error: %s\n", traceID, appErr)

	var statusCode int
	switch appErr.ErrCode {
	case NAData:
		statusCode = http.StatusFound
	case NoTargetData, ReqBodyDecodeFailed, BadParam:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(appErr)
}
