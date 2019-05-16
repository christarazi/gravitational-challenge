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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/christarazi/gravitational-challenge/client/api"
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
	id, err := api.NewClient(args).Start()
	if err != nil {
		return err
	}

	fmt.Printf("%d\n", id)

	return nil
}

func init() {
	rootCmd.AddCommand(startCmd)
}
