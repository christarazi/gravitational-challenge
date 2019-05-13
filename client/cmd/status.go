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
	"net/http"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/christarazi/gravitational-challenge/client/util"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [Job ID]",
	Short: "Retrieve the status of all jobs or a single job based on ID",
	Long: `This command supports taking in a Job ID which will retrieve the
status of that job, or when given no arguments, it will return all the
jobs.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			doAllStatus()
			return nil
		}

		id, err := util.ConvertAndValidateID(args[0])
		if err != nil {
			return err
		}
		doStatus(id)
		return nil
	},
}

// This is the response data structure when the user does not supply a Job ID.
// All the jobs will be returned.
type allStatusResponse struct {
	Jobs []struct {
		ID      uint64   `json:"id"`
		Command string   `json:"command"`
		Args    []string `json:"args"`
		Status  string   `json:"status"`
	} `json:"jobs"`
}

// This is the response data structure when the user supplies a Job ID. The
// single status of the job will be returned.
type statusResponse struct {
	Status string `json:"status"`
}

func do(uri string) *http.Response {
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting response: %v", err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading body of response: %v", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Server returned %d: %v", resp.StatusCode, string(body))
		os.Exit(1)
	}

	return resp
}

func doAllStatus() {
	uri := fmt.Sprintf("http://0.0.0.0:8080/status")

	resp := do(uri)
	defer resp.Body.Close()

	asr := &allStatusResponse{}
	err := json.NewDecoder(resp.Body).Decode(asr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding response: %v", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Command", "Args", "Status"})

	for _, v := range asr.Jobs {
		str := fmt.Sprintf("%d|%s|%v|%s", v.ID, v.Command, v.Args, v.Status)
		table.Append(strings.Split(str, "|"))
	}

	table.Render()
}

func doStatus(id uint64) {
	uri := fmt.Sprintf("http://0.0.0.0:8080/status/%d", id)

	resp := do(uri)
	defer resp.Body.Close()

	sr := &statusResponse{}
	err := json.NewDecoder(resp.Body).Decode(sr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding response: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", sr.Status)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
