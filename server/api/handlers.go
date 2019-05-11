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

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func validateID(id uint64) error {
	// TODO: This needs a mutex if accessing global Jobs
	if (id - 1) >= uint64(len(models.Jobs)) {
		return fmt.Errorf("job with id %v does not exit", id)
	}

	return nil
}

func GetJobStatus(w http.ResponseWriter, r *http.Request) {
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

	err = validateID(id)
	if err != nil {
		msg := fmt.Sprintf("/status error: %v", err)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// TODO: Mutex here.
	j := models.Jobs[id-1]

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

func StartJob(w http.ResponseWriter, r *http.Request) {
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

	// TODO: This will need a mutex around // it.
	models.Jobs = append(models.Jobs, j)
	j.ID = uint64(len(models.Jobs))

	log.Printf("/start: created new job with id %d", j.ID)

	j.Status = "Running"
	// TODO: Actually start running the job here.
	log.Printf("/start: running job with id %d", j.ID)

	// TODO: Move this out into a separate file.
	type startResponse struct {
		ID uint64 `json:"id"`
	}

	err = json.NewEncoder(w).Encode(startResponse{ID: j.ID})
	if err != nil {
		msg := fmt.Sprintf("/start error: %v", err)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}

func StopJob(w http.ResponseWriter, r *http.Request) {
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
	err = validateID(id)
	if err != nil {
		msg := fmt.Sprintf("/stop error: %v", err)
		log.Println(msg)

		// TODO: Return API specific error codes. For example, if no jobs
		// exist, it would be 4xx.
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// TODO: // Actually stop running the job here.
	log.Printf("/stop: stopped job with id %d", id)

	// TODO: // Need a mutex here.
	j := models.Jobs[id-1]

	// TODO: // Actually get the exit code here.
	j.Status = fmt.Sprintf("Stopped (ec: %d)", 42)

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
