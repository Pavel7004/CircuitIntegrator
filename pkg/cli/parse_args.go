package cli

import (
	"github.com/pborman/getopt/v2"
)

var (
	Filename string  = "result.png"
	Step     float64 = 0.1
	Dpi      int     = 40
	LoadRes  float64 = 9000.0
)

func init() {
	getopt.FlagLong(&Filename, "output", 'o', "Output file path")
	getopt.FlagLong(&Step, "step", 's', "Integrator step")
	getopt.FlagLong(&Dpi, "dpi", 'd', "Plot dpi")
	getopt.FlagLong(&LoadRes, "load-resistance", 'l', "Set load resistance")
}

func ParseArgs() {
	getopt.Parse()
}
