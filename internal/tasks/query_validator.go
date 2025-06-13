package tasks

import (
	"log"
	"net/http"
	"strconv"
)

const (
	Ordered   = "ordered"
	Created   = "created"
	Running   = "running"
	Completed = "completed"
	Failed    = "failed"
)

var params = []string{Created, Running, Completed, Failed}

func isOrdered(r *http.Request) bool {
	ordered := r.URL.Query().Get(Ordered)
	if ordered != "" {
		order, parseErr := strconv.ParseBool(ordered)
		if parseErr != nil {
			log.Printf("wrong \"%s\" parameter value: %s", Ordered, ordered)
			return false
		}
		return order
	}
	return false
}

func getValidQueryParams(r *http.Request) map[string]bool {
	validQueryParams := map[string]bool{}

	for _, param := range params {
		queryParam := r.URL.Query().Get(param)
		if queryParam != "" {
			parsedParam, parseErr := strconv.ParseBool(queryParam)
			if parseErr != nil {
				log.Printf("wrong \"%s\" parameter value: %s", param, queryParam)
				continue
			}
			validQueryParams[param] = parsedParam
		}
	}
	return validQueryParams
}
