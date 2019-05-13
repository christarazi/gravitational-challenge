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
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/christarazi/gravitational-challenge/server/app"
)

func main() {
	app := app.App{}
	app.Initialize()

	// Set up signal handlers for the following signals for graceful shutdown.
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, []os.Signal{
		os.Interrupt,
		syscall.SIGABRT,
		syscall.SIGQUIT,
		syscall.SIGTERM}...)

	// TODO: This is hard coded for now. In the future, we can have a
	// configurable address / port number.
	port := "8080"
	server := &http.Server{Addr: ":" + port, Handler: app.Router}

	go func() {
		log.Printf("Listening on http://0.0.0.0:%s\n", port)

		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	<-stopCh

	log.Println("Shutting down the server...")
	server.Shutdown(context.Background())
	log.Println("Server gracefully shutdown")
}
