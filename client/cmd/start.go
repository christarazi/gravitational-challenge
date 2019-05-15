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
	"net/http"

	"github.com/spf13/cobra"

	"github.com/christarazi/gravitational-challenge/client/util"
	"github.com/christarazi/gravitational-challenge/config"
	"github.com/christarazi/gravitational-challenge/models"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start -- <command> <args...>",
	Short: "Start a job on the server",
	Long: `This command requires a '--' after the 'start' in order to parse the
command and arguments correctly.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return doStart(args)
	},
}

func doStart(args []string) error {
	data, err := json.Marshal(models.StartRequest{
		Command: args[0],
		Args:    args[1:],
	})
	if err != nil {
		return fmt.Errorf("Error marshalling request: %v", err)
	}

	uri := fmt.Sprintf("http://0.0.0.0:%d/start", config.Port)
	resp, err := http.Post(uri, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("Error getting response: %v", err)
	}

	defer resp.Body.Close()

	err = util.CheckHTTPStatusCode(resp)
	if err != nil {
		return err
	}

	sr := &models.StartResponse{}
	err = json.NewDecoder(resp.Body).Decode(sr)
	if err != nil {
		return fmt.Errorf("Error decoding response: %v", err)
	}

	fmt.Printf("%d\n", sr.ID)

	return nil
}

func init() {
	rootCmd.AddCommand(startCmd)
}
