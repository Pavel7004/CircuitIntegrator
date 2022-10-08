package midpoint

import (
	"context"
	"math"

	"github.com/Pavel7004/Common/tracing"

	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator"
	"github.com/Pavel7004/GraphPlot/pkg/circuit"
)

type MidpointInt struct {
	begin  float64
	end    float64
	step   float64
	saveFn func(t float64, x *circuit.Circuit)
}

var _ integrator.Integrator = (*MidpointInt)(nil)

func NewMidpointInt(begin, end, step float64, saveFn func(t float64, x *circuit.Circuit)) integrator.Integrator {
	return &MidpointInt{
		begin:  begin,
		end:    end,
		step:   step,
		saveFn: saveFn,
	}
}

func (mi *MidpointInt) Integrate(ctx context.Context, circ *circuit.Circuit) float64 {
	span, _ := tracing.StartSpanFromContext(ctx)
	span.SetTag("StartPoint", mi.begin)
	span.SetTag("EndPoint", mi.end)
	span.SetTag("Step", mi.step)
	span.SetTag("RK-stages", 2)

	defer span.Finish()

	var (
		t    = mi.begin
		last bool
	)

	for !last && mi.step > 0 {
		if t+mi.step > mi.end {
			last = true
			mi.step = mi.end - t
		}

		k1 := circ.GetDerivative()
		k2 := circ.Clone().ApplyDerivative(mi.step/2, k1).GetDerivative()

		if !circ.CheckDerivative(mi.step, k2) {
			mi.step = circ.CalculateOptimalStep(mi.step, k2)
			mi.step = math.Cbrt(mi.step)

			last = true
		}

		circ.ApplyDerivative(mi.step, k2)
		t += mi.step

		mi.saveFn(t, circ)
	}

	span.SetTag("finish-point", t)
	span.SetTag("finish-step", mi.step)

	return t
}
