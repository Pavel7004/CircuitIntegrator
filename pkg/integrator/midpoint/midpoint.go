package midpoint

import (
	"context"

	"github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/common/tracing"
	"github.com/Pavel7004/GraphPlot/pkg/integrator"
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

func (mi *MidpointInt) Integrate(ctx context.Context, circ *circuit.Circuit) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	span.SetTag("StartPoint", mi.begin)
	span.SetTag("EndPoint", mi.end)
	span.SetTag("Step", mi.step)
	span.SetTag("RK-stages", 2)

	defer span.Finish()

	var (
		t    = mi.begin
		last bool
	)

	for !last {
		if t+mi.step > mi.end {
			last = true
			mi.step = mi.end - t
		}

		k1 := circ.Clone()
		k1.ApplyDerivative(mi.step/2, k1.GetDerivative())
		circ.ApplyDerivative(mi.step, k1.GetDerivative())
		t += mi.step
		mi.saveFn(t, circ)
	}
}
