package webserver

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(router chi.Router,
	handlers map[string]http.HandlerFunc,
	webServerPort string) *WebServer {
	return &WebServer{Router: router,
		Handlers:      handlers,
		WebServerPort: webServerPort}
}

func (w *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	w.Handlers[path] = handler
}

func (w *WebServer) Start() {
	w.Router.Use(middleware.Logger)
	for path, handler := range w.Handlers {
		w.Router.Handle(path, handler)
	}
	http.ListenAndServe(w.WebServerPort, w.Router)
}
