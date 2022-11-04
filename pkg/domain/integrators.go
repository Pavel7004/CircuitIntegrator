package domain

import (
	misc "github.com/Pavel7004/Common/misc"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator"
	bogatskiyshampin "github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/bogatskiy-Shampin"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/euler"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/midpoint"
	midpointimplicit "github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/midpoint-implicit"
	threeeighth "github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/three-eighth"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/trapeziod"
)

var (
	Integrators = []integrator.NewIntFunc{
		euler.NewEulerInt,
		midpoint.NewMidpointInt,
		midpointimplicit.NewMidpointImplInt,
		bogatskiyshampin.NewShampinInt,
		threeeighth.NewThreeEighthInt,
		trapeziod.NewTrapezoidInt,
	}

	IntegratorsNames []string
)

func init() {
	IntegratorsNames = make([]string, 0, len(Integrators))
	for _, int := range Integrators {
		IntegratorsNames = append(IntegratorsNames, misc.GetFuncModule(int))
	}
}
