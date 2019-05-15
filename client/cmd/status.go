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
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/christarazi/gravitational-challenge/client/util"
	"github.com/christarazi/gravitational-challenge/config"
	"github.com/christarazi/gravitational-challenge/models"
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
			return doAllStatus()
		}

		id, err := util.ConvertAndValidateID(args[0])
		if err != nil {
			return err
		}
		return doStatus(id)
	},
}

func doAllStatus() error {
	uri := fmt.Sprintf("http://0.0.0.0:%d/status", config.Port)

	resp, err := do(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	asr := &models.AllStatusResponse{}
	err = json.NewDecoder(resp.Body).Decode(asr)
	if err != nil {
		return fmt.Errorf("Error decoding response: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Command", "Args", "Status"})

	for _, v := range asr.Jobs {
		table.Append([]string{
			strconv.FormatUint(v.ID, 10),
			v.Command,
			strings.Join(v.Args, ","),
			v.Status})
	}

	table.Render()

	return nil
}

func doStatus(id uint64) error {
	uri := fmt.Sprintf("http://0.0.0.0:8080/status/%d", id)

	resp, err := do(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sr := &models.StatusResponse{}
	err = json.NewDecoder(resp.Body).Decode(sr)
	if err != nil {
		return fmt.Errorf("Error decoding response: %v", err)
	}

	fmt.Printf("%s\n", sr.Status)

	return nil
}

func do(uri string) (*http.Response, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("Error getting response: %v", err)
	}

	err = util.CheckHTTPStatusCode(resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
