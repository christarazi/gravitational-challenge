/*
Copyright © 2019 Chris Tarazi <tarazichris@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/christarazi/gravitational-challenge/models"
	"github.com/christarazi/gravitational-challenge/server/manager"

	"github.com/gorilla/mux"
)

// GetAllJobStatus implements the /status endpoint which returns all the jobs
// (models.Job).
func GetAllJobStatus(m *manager.Manager, w http.ResponseWriter, r *http.Request) {
	statusResponse := models.AllStatusResponse{Jobs: m.Jobs()}
	err := json.NewEncoder(w).Encode(statusResponse)
	if err != nil {
		reportHTTPError(&w, fmt.Sprintf("/status error: %v", err),
			http.StatusBadRequest)
		return
	}
}

// GetJobStatus implements the /status endpoint that takes an ID of a job to
// return the status of.
func GetJobStatus(m *manager.Manager, w http.ResponseWriter, r *http.Request) {
	j, err := m.JobStatus(mux.Vars(r)["id"])
	if err != nil {
		reportHTTPError(&w, fmt.Sprintf("/status error: %v", err),
			http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(j)
	if err != nil {
		reportHTTPError(&w, fmt.Sprintf("/status error: %v", err),
			http.StatusBadRequest)
		return
	}
}

// StartJob implements the /start endpoint that takes the command and the args
// of a job to start from the request body.
func StartJob(m *manager.Manager, w http.ResponseWriter, r *http.Request) {
	j := &models.Job{}

	err := json.NewDecoder(r.Body).Decode(j)
	if err != nil {
		reportHTTPError(&w, fmt.Sprintf("/start error: %v", err),
			http.StatusBadRequest)
		return
	}

	id, err := m.StartJob(j)
	if err != nil {
		reportHTTPError(&w, fmt.Sprintf("/start error: %v", err),
			http.StatusBadRequest)
		return
	}

	log.Printf("/start: running job with id %d", id)

	err = json.NewEncoder(w).Encode(models.StartResponse{ID: id})
	if err != nil {
		reportHTTPError(&w, fmt.Sprintf("/start error: %v", err),
			http.StatusBadRequest)
		return
	}
}

// StopJob implements the /stop endpoint which stops a job based on ID.
func StopJob(m *manager.Manager, w http.ResponseWriter, r *http.Request) {
	request := &models.StopRequest{}

	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		reportHTTPError(&w, fmt.Sprintf("/stop error: %v", err),
			http.StatusBadRequest)
		return
	}

	id := request.ID
	err = m.StopJob(id)
	if err != nil {
		reportHTTPError(&w, fmt.Sprintf("/stop error: %v", err),
			http.StatusBadRequest)
		return
	}

	log.Printf("/stop: job %d stopped\n", id)

	err = json.NewEncoder(w).Encode([]byte{})
	if err != nil {
		reportHTTPError(&w, fmt.Sprintf("/stop error: %v", err),
			http.StatusBadRequest)
		return
	}
}

func reportHTTPError(w *http.ResponseWriter, msg string, statusCode int) {
	log.Println(msg)

	// TODO: Return API specific error codes. For example, if no jobs
	// exist, it would be 4xx.
	http.Error(*w, msg, http.StatusBadRequest)
}
