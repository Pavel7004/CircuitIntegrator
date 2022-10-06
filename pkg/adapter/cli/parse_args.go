package cli

import (
	"github.com/pborman/getopt/v2"
)

var (
	DirName       string  = "results"
	Step          float64 = 0.1
	Dpi           int     = 40
	LoadRes       float64 = 10000.0
	CapCount      uint    = 5
	PointBuffSize int     = 100
)

func init() {
	getopt.FlagLong(&DirName, "output", 'o', "Output file path")
	getopt.FlagLong(&Step, "step", 's', "Integrator step")
	getopt.FlagLong(&Dpi, "dpi", 'd', "Plot dpi")
	getopt.FlagLong(&LoadRes, "load-resistance", 'l', "Set load resistance")
	getopt.FlagLong(&CapCount, "stages-count", 'c', "Set number of capacitors stages in circuit")
	getopt.FlagLong(&PointBuffSize, "buffer-size", 'b', "Set size of points buffer to optimize plotting speed")
}

func ParseArgs() {
	getopt.Parse()
}
