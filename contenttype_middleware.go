package usrsvc

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func acceptContentType() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return handlers.ContentTypeHandler(h, []string{"application/json"}...)
	}
}

func addContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
