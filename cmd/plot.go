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
	"github.com/spf13/cobra"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
	plotcli "github.com/Pavel7004/GraphPlot/pkg/components/plot-cli"
)

var (
	supplyVol *float64
	capCount  *uint
	loadRes   *float64
	step      *float64
	output    *string
	buffSize  *int
	dpi       *int
)

var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Generate plot image",
	Long: `Generate plot image in your directory.

Due to implementation don't support step less than 0.01.

Example: graph plot -d 300 -s 0.1 -o results

This will create directory results/ and put plot images into it.`,
	Run: func(cmd *cobra.Command, args []string) {
		chargeCirc := circuit.ChargeComponents{
			SupplyVoltage:     *supplyVol,
			Capacity:          0.001,
			Resistance:        5000,
			StagesCount:       *capCount,
			GapTriggerVoltage: 5700,
			HoldingVoltage:    1,
		}

		load := circuit.LoadComponents{
			Resistance: *loadRes,
		}

		circ := circuit.New(chargeCirc, load)
		p := plotcli.NewPlotterCli(circ, &plotcli.Settings{
			Step:       *step,
			FolderName: *output,
			BuffSize:   *buffSize,
			Dpi:        *dpi,
		})

		p.Plot()
	},
}

func init() {
	rootCmd.AddCommand(plotCmd)

	capCount = plotCmd.Flags().UintP("capacitors", "c", 6, "Change number of capacitors in circuit")
	supplyVol = plotCmd.Flags().Float64P("supply-voltage", "v", 6000, "Change supply voltage in circuit")
	loadRes = plotCmd.Flags().Float64P("load-resistance", "l", 10000, "Change load resistance value")

	step = plotCmd.Flags().Float64P("step", "s", 0.1, "Change default step amount")
	output = plotCmd.Flags().StringP("output", "o", "results", "Change results directory name")
	buffSize = plotCmd.Flags().Int("buffer-size", 100, "Change size of line-draw buffer")
	dpi = plotCmd.Flags().Int("dpi", 320, "Change dpi of resulting images")
}
