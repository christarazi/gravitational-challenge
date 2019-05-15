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

package models

import (
	"os/exec"
)

// Job is the main data structure of the API.
type Job struct {
	ID      uint64   `json:"id"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Status  string   `json:"status"`

	Process *exec.Cmd
}

// AllStatusResponse is the JSON structure for the /status endpoint when the
// user does not supply a Job ID.  All the jobs will be returned.
type AllStatusResponse struct {
	Jobs []struct {
		ID      uint64   `json:"id"`
		Command string   `json:"command"`
		Args    []string `json:"args"`
		Status  string   `json:"status"`
	} `json:"jobs"`

	// Jobs []Job `json:"jobs"`
}

// StatusResponse is the JSON response structure for the /status endpoint when
// the user supplies a Job ID. The single status of the job will be returned.
type StatusResponse struct {
	Status string `json:"status"`
}

// StartRequest is the JSON request structure for the /start endpoint.
type StartRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// StartResponse is the JSON response structure for the /start endpoint.
type StartResponse struct {
	ID uint64 `json:"id"`
}

// StopRequest is the JSON response structure for the /stop endpoint.
type StopRequest struct {
	ID uint64 `json:"id"`
}
