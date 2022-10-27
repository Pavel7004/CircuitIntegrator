package pointgenerator

import (
	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator"
	"github.com/Pavel7004/GraphPlot/pkg/circuit"
)

type Args struct {
	Circuit *circuit.Circuit
	Step    float64

	SaveFn   func(t float64, x *circuit.Circuit)
	NewIntFn integrator.NewIntFunc
}
