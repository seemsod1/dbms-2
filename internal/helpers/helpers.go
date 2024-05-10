package helpers

import (
	"fmt"
	"net/http"
)

// ServerError sends a server error response
func ServerError(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusInternalServerError)
}
