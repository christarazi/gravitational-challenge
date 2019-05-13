package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/christarazi/gravitational-challenge/server/app/handler"
	"github.com/christarazi/gravitational-challenge/server/manager"
)

type App struct {
	Router  *mux.Router
	Manager *manager.Manager
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.Manager = manager.NewManager()

	a.setRoutes()
}

// setRoutes sets the all required routers
func (a *App) setRoutes() {
	a.Get("/status", a.GetAllJobStatus)
	a.Get("/status/{id:[0-9]+}", a.GetJobStatus)
	a.Post("/start", a.StartJob)
	a.Post("/stop", a.StopJob)
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) GetAllJobStatus(w http.ResponseWriter, r *http.Request) {
	handler.GetAllJobStatus(a.Manager, w, r)
}

func (a *App) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	handler.GetJobStatus(a.Manager, w, r)
}

func (a *App) StartJob(w http.ResponseWriter, r *http.Request) {
	handler.StartJob(a.Manager, w, r)
}

func (a *App) StopJob(w http.ResponseWriter, r *http.Request) {
	handler.StopJob(a.Manager, w, r)
}
