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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [Job ID]",
	Short: "Retrieve the status of a job",
	Long: `This command supports taking in a Job ID which will retrieve the
status of that job, or when given no arguments, it will return all the
jobs.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Validate the arg is a number.
		doStatus(args)
	},
}

type statusResponse struct {
	Status string `json:"status"`
}

func doStatus(args []string) {
	uri := fmt.Sprintf("http://0.0.0.0:8080/status/%v", args[0])
	resp, err := http.Get(uri)
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

	sr := &statusResponse{}
	err = json.NewDecoder(resp.Body).Decode(sr)
	if err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	fmt.Printf("%s\n", sr.Status)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
