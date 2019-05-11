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

// Job is the main data structure that API will utilize.
type Job struct {
	ID      uint64   `json:"id"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	// Stdout  []byte   `json:"stdout"`
	// Stderr  []byte   `json:"stderr"`
	Status string `json:"status"`
}

// Jobs holds all the processes that were requested to start by the client.
var Jobs []*Job

func NewJob(id uint64, command string, args []string) *Job {
	return &Job{
		ID:      id,
		Command: command,
		Args:    args,
		Status:  "",
	}
}
