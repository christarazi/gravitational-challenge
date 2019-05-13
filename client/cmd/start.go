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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start -- <command> <args...>",
	Short: "Start a job on the server",
	Long: `This command requires a '--' after the 'start' in order to parse the
command and arguments correctly.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		doStart(args)
	},
}

type startRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type startResponse struct {
	ID uint64 `json:"id"`
}

func doStart(args []string) {
	data, err := json.Marshal(startRequest{
		Command: args[0],
		Args:    args[1:],
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling request: %v", err)
		os.Exit(1)
	}

	resp, err := http.Post("http://0.0.0.0:8080/start",
		"application/json", bytes.NewReader(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting response: %v", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading body of response: %v", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Server returned %d: %v", resp.StatusCode, string(body))
		os.Exit(1)
	}

	sr := &startResponse{}
	err = json.NewDecoder(resp.Body).Decode(sr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding response: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%d\n", sr.ID)
}

func init() {
	rootCmd.AddCommand(startCmd)
}
