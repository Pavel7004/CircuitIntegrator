/*
Copyright © 2022 Kovalev Pavel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/http"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long: `Start server hosting plotting website.

Example: graph server

This will start server on localhost:8088. To modify hostname and port create config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		s := http.New()

		if err := s.Run(); err != nil {
			log.Error().Err(err).Msg("Error during server operation.")
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
