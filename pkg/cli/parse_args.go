package cli

import (
	"os"

	"github.com/pborman/getopt/v2"
)

var (
	fileName       string = os.Getenv("PWD") + string(os.PathSeparator) + "result.png"
	capacitorIndex uint   = 1
	state          string = "help"
)

func init() {
	getopt.FlagLong(&fileName, "outputFile", 'o', "output file path")
	getopt.Flag(&capacitorIndex, 'c', "Capacitor Index")
}

func ParseArgs() {
	getopt.Parse()
}
