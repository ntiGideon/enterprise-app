package helpers

import (
	"encoding/json"
	"net/http"
)

func ReadRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	if err != nil {
		PanicAllErrors(err)
	}
}

func WriteResponseBody(w http.ResponseWriter, result interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		PanicAllErrors(err)
	}
}
