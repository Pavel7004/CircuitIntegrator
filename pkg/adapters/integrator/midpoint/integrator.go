package midpoint

import (
	"context"

	"github.com/Pavel7004/Common/tracing"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator"
)

type MidpointInt struct {
	begin  float64
	end    float64
	step   float64
	saveFn func(t float64, x *circuit.Circuit) error
}

var _ integrator.Integrator = (*MidpointInt)(nil)

func NewMidpointInt(begin, end, step float64, saveFn func(t float64, x *circuit.Circuit) error) integrator.Integrator {
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
		t = mi.begin

		last bool
	)

	for !last {
		if t+mi.step > mi.end {
			last = true
			mi.step = mi.end - t
		}

		k1 := circ.GetDerivative()
		k2 := circ.Clone().ApplyDerivative(mi.step/2, k1).GetDerivative()

		if !circ.CheckDerivative(mi.step, k2) {
			mi.step = circ.CalculateOptimalStep(mi.step, k2)

			last = true
		}

		circ.ApplyDerivative(mi.step, k2)
		t += mi.step

		if err := mi.saveFn(t, circ); err != nil {
			span.SetTag("Error", err)
			break
		}
	}

	span.SetTag("finish-point", t)
	span.SetTag("finish-step", mi.step)

	return t
}
