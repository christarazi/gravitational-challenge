package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/christarazi/gravitational-challenge/server/app/handler"
	"github.com/christarazi/gravitational-challenge/server/manager"
)

// App struct is the main data structure for the server.
type App struct {
	Router  *mux.Router
	Manager *manager.Manager
}

// Initialize creates a new instance of App and sets the routes for the server.
func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.Manager = manager.NewManager()

	a.setRoutes()
}

func (a *App) setRoutes() {
	a.Get("/status", a.GetAllJobStatus)
	a.Get("/status/{id:[0-9]+}", a.GetJobStatus)
	a.Post("/start", a.StartJob)
	a.Post("/stop", a.StopJob)
}

// Get wraps the router for all the GET endpoints. This allows us to inject the
// Manager into the handler.
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for all the POST endpoints. This allows us to inject
// the Manager into the handler.
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// GetAllJobStatus forwards the request to the endpoint implementation in the
// handler package.
func (a *App) GetAllJobStatus(w http.ResponseWriter, r *http.Request) {
	handler.GetAllJobStatus(a.Manager, w, r)
}

// GetJobStatus forwards the request to the endpoint implementation in the
// handler package.
func (a *App) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	handler.GetJobStatus(a.Manager, w, r)
}

// StartJob forwards the request to the endpoint implementation in the handler
// package.
func (a *App) StartJob(w http.ResponseWriter, r *http.Request) {
	handler.StartJob(a.Manager, w, r)
}

// StopJob forwards the request to the endpoint implementation in the handler
// package.
func (a *App) StopJob(w http.ResponseWriter, r *http.Request) {
	handler.StopJob(a.Manager, w, r)
}
