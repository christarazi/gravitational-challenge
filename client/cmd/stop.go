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

package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop <ID>",
	Short: "Stop a job on the server",
	Long:  `This command stops a job on the server based on ID`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Validate the args here.
		doStop(args)
	},
}

type stopRequest struct {
	ID uint64 `json:"id"`
}

func doStop(args []string) {
	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		log.Fatalf("Error converting '%v' to integer: %v", args[0], err)
	}

	data, err := json.Marshal(stopRequest{
		ID: id,
	})

	resp, err := http.Post("http://0.0.0.0:8080/stop",
		"application/json", bytes.NewReader(data))
	if err != nil {
		log.Fatalf("Error getting response: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading body of response: %v", err)
		}

		log.Fatalf("Server returned %d: %v", resp.StatusCode, string(body))
	}
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
