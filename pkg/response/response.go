package response

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(writer http.ResponseWriter, data interface{}, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	if statusCode == http.StatusNoContent {
		return
	}
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		panic(err)
	}
}
