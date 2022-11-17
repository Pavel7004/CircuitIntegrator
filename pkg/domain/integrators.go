package domain

import (
	misc "github.com/Pavel7004/Common/misc"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator"
	bogatskiyshampin "github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/bogatskiy-Shampin"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/euler"
	eulerimplicit "github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/euler-implicit"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/midpoint"
	threeeighth "github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/three-eighth"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/trapeziod"
)

var (
	Integrators = []integrator.NewIntFunc{
		euler.NewEulerInt,
		midpoint.NewMidpointInt,
		bogatskiyshampin.NewShampinInt,
		threeeighth.NewThreeEighthInt,
		eulerimplicit.NewEulerImplInt,
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
