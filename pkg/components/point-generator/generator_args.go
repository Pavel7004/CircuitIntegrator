package pointgenerator

import (
	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator"
)

type Args struct {
	Circuit *circuit.Circuit
	Step    float64

	SaveFn   func(t float64, x *circuit.Circuit)
	NewIntFn integrator.NewIntFunc
}
