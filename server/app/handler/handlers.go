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
	"strconv"

	"github.com/christarazi/gravitational-challenge/server/manager"
	"github.com/christarazi/gravitational-challenge/server/models"

	"github.com/gorilla/mux"
)

func convertIDToUint(str string) (uint64, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetAllJobStatus(m *manager.Manager, w http.ResponseWriter, r *http.Request) {
	// TODO: Move this out into a separate file.
	statusResponse := struct {
		Jobs []*models.Job `json:"jobs"`
	}{Jobs: m.GetJobs()}

	err := json.NewEncoder(w).Encode(statusResponse)

	if err != nil {
		msg := fmt.Sprintf("/status error: %v", err)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}

func GetJobStatus(m *manager.Manager, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("/status vars: %v\n", vars)

	id, err := convertIDToUint(vars["id"])
	if err != nil {
		msg := fmt.Sprintf("/status error: %v", err)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if !m.IsAJob(id) {
		msg := fmt.Sprintf("/status error: job with id %v does not exist", id)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	j := m.GetJobByID(id)

	err = json.NewEncoder(w).Encode(j)
	if err != nil {
		msg := fmt.Sprintf("/status error: %v", err)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}

func StartJob(m *manager.Manager, w http.ResponseWriter, r *http.Request) {
	j := &models.Job{}

	err := json.NewDecoder(r.Body).Decode(j)
	if err != nil {
		msg := fmt.Sprintf("/start error: %v", err)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id, err := m.AddAndStartJob(j)
	if err != nil {
		msg := fmt.Sprintf("/start failed to start job %d: %v", id, err)
		log.Println(msg)

		m.SetJobStatus(j, "Errored")

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	log.Printf("/start: running job with id %d", id)

	// TODO: Move this out into a separate file.
	type startResponse struct {
		ID uint64 `json:"id"`
	}

	err = json.NewEncoder(w).Encode(startResponse{ID: id})
	if err != nil {
		msg := fmt.Sprintf("/start error: %v", err)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}

func StopJob(m *manager.Manager, w http.ResponseWriter, r *http.Request) {
	// TODO: Move this into separate file.
	type stopRequest struct {
		ID uint64 `json:"id"`
	}

	request := &stopRequest{}

	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		msg := fmt.Sprintf("/stop error: %v", err)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id := request.ID
	if !m.IsAJob(id) {
		msg := fmt.Sprintf("/stop error: job id %d does not exist", id)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = m.StopJobByID(id)
	if err != nil {
		msg := fmt.Sprintf("/stop error: %v", err)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	log.Printf("/stop: job %d stopped\n", id)

	err = json.NewEncoder(w).Encode([]byte{})
	if err != nil {
		msg := fmt.Sprintf("/stop error: %v", err)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}
