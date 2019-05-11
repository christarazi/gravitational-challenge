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

package main

import (
	"log"
	"net/http"

	"github.com/christarazi/gravitational-challenge/server/api"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/status/{id:[0-9]+}", api.GetJobStatus).Methods("GET")
	router.HandleFunc("/start", api.StartJob).Methods("POST")
	router.HandleFunc("/stop", api.StopJob).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
