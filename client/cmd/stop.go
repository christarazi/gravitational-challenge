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

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop <ID>",
	Short: "Stop a job on the server",
	Long:  `This command stops a job on the server based on ID`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := util.ConvertAndValidateID(args[0])
		if err != nil {
			return err
		}
		return doStop(id)
	},
}

func doStop(id uint64) error {
	data, err := json.Marshal(models.StopRequest{ID: id})
	if err != nil {
		return fmt.Errorf("Error marshalling request: %v", err)
	}

	uri := fmt.Sprintf("http://0.0.0.0:%d/stop", config.Port)
	resp, err := http.Post(uri, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("Error getting response: %v", err)
	}

	defer resp.Body.Close()

	err = util.CheckHTTPStatusCode(resp)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
