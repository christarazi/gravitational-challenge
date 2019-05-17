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

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [Job ID]",
	Short: "Retrieve the status of all jobs or a single job based on ID",
	Long: `This command supports taking in a Job ID which will retrieve the
status of that job, or when given no arguments, it will return all the
jobs.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := api.NewClient(args).Status()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", s)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
