package router

import "net/http"

type Registrar interface {
	RegisterRoutes(mux *http.ServeMux)
}

type Middleware = func(http.Handler) http.Handler

func New(mws []Middleware, registrars ...Registrar) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	for _, r := range registrars {
		r.RegisterRoutes(mux)
	}

	var h http.Handler = mux
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}