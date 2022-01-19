package cli

import (
	"github.com/pborman/getopt/v2"
)

var (
	Filename string  = "result.png"
	Step     float64 = 0.1
	Dpi      int     = 40
	LoadRes  float64 = 10000.0
	CapCount uint    = 5
)

func init() {
	getopt.FlagLong(&Filename, "output", 'o', "Output file path")
	getopt.FlagLong(&Step, "step", 's', "Integrator step")
	getopt.FlagLong(&Dpi, "dpi", 'd', "Plot dpi")
	getopt.FlagLong(&LoadRes, "load-resistance", 'l', "Set load resistance")
	getopt.FlagLong(&CapCount, "stages-count", 'c', "Set number of capacitors stages in circuit")
}

func ParseArgs() {
	getopt.Parse()
}
