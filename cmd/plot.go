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
	supplyVol  float64
	capCount   uint
	capacity   float64
	resistance float64
	triggerVol float64
	loadRes    float64
	step       float64
	output     string
	buffSize   int
	format     string
)

var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Generate plot image",
	Long: `Generate plot images in your directory.
Available formats are: png, svg, tex, pdf, jpg, jpeg, eps, tif, tiff.

Example: graph plot -s 0.1 -o results -f png

This will create directory results/ and put plot images in "png" format into it.`,
	Run: func(cmd *cobra.Command, args []string) {
		chargeCirc := &circuit.ChargeComponents{
			SupplyVoltage:     supplyVol,
			Capacity:          capacity / 1000,
			Resistance:        resistance,
			StagesCount:       capCount,
			GapTriggerVoltage: triggerVol,
			HoldingVoltage:    1,
		}

		load := &circuit.LoadComponents{
			Resistance: loadRes,
		}

		circ := circuit.New(chargeCirc, load)
		p := plotcli.NewPlotterCli(circ, &plotcli.Settings{
			Step:       step,
			FolderName: output,
			Format:     format,
			BuffSize:   buffSize,
		})

		p.Plot()
	},
}

func init() {
	rootCmd.AddCommand(plotCmd)

	plotCmd.Flags().UintVarP(&capCount, "capacitors", "n", 6, "change number of capacitors in circuit")
	plotCmd.Flags().Float64VarP(&supplyVol, "supply-voltage", "v", 6000, "change supply voltage in circuit")
	plotCmd.Flags().Float64VarP(&loadRes, "load-resistance", "l", 10000, "change load resistance value")
	plotCmd.Flags().Float64VarP(&capacity, "capacity", "c", 1, "change capacity of circuit (in mF)")
	plotCmd.Flags().Float64VarP(&resistance, "resistance", "r", 5000, "change resistance of circuit")
	plotCmd.Flags().Float64VarP(&triggerVol, "trigger-voltage", "t", 5700, "change gap trigger voltage in circuit")

	plotCmd.Flags().Float64VarP(&step, "step", "s", 0.001, "change default step amount")
	plotCmd.Flags().StringVarP(&output, "output", "o", "results", "change results directory name")
	plotCmd.Flags().StringVarP(&format, "format", "f", "svg", "change resulting images format")

	plotCmd.Flags().IntVar(&buffSize, "buffer-size", 100, "change size of line-draw buffer")
}
