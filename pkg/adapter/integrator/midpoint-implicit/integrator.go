package midpointimplicit

import (
	"context"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator"
	"github.com/Pavel7004/GraphPlot/pkg/circuit"
)

type MidpointImpInt struct {
	begin  float64
	end    float64
	step   float64
	saveFn func(t float64, x *circuit.Circuit)
}

var _ integrator.Integrator = (*MidpointImpInt)(nil)

func NewMidpointImplInt(begin, end, step float64, saveFn func(t float64, x *circuit.Circuit)) integrator.Integrator {
	return &MidpointImpInt{
		begin:  begin,
		end:    end,
		step:   step,
		saveFn: saveFn,
	}
}

func (si *MidpointImpInt) Integrate(ctx context.Context, circ *circuit.Circuit) float64 {
	span, ctx := tracing.StartSpanFromContext(ctx)
	span.SetTag("StartPoint", si.begin)
	span.SetTag("EndPoint", si.end)
	span.SetTag("Step", si.step)
	span.SetTag("RK-stages", 2)

	defer span.Finish()

	var (
		t    = si.begin
		last bool
	)

	for !last && si.step > 0 {
		if t+si.step > si.end {
			last = true
			si.step = si.end - t
		}

		k1 := circ.GetDerivative()

		if !circ.CheckDerivative(si.step, k1) {
			si.step = circ.CalculateOptimalStep(si.step, k1)

			last = true
		}

		circ.ApplyDerivative(si.step, k1)
		circ.ApplyDerivative(si.step*si.step/2, k1)
		t += si.step

		si.saveFn(t, circ)
	}

	span.SetTag("finish-point", t)
	span.SetTag("finish-step", si.step)

	return t
}
