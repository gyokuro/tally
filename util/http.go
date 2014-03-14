package util

import (
	"net/http"
)

func AddHeaders(w *http.ResponseWriter, headers map[string]string) {
	for key, value := range headers {
		(*w).Header().Add(key, value)
	}
}
