package trapeziod

import (
	"context"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/integrator"
)

type TrapezoidInt struct {
	begin  float64
	end    float64
	step   float64
	saveFn func(t float64, x *circuit.Circuit)
}

var _ integrator.Integrator = (*TrapezoidInt)(nil)

func NewTrapezoidInt(begin, end, step float64, saveFn func(t float64, x *circuit.Circuit)) integrator.Integrator {
	return &TrapezoidInt{
		begin:  begin,
		end:    end,
		step:   step,
		saveFn: saveFn,
	}
}

func (si *TrapezoidInt) Integrate(ctx context.Context, circ *circuit.Circuit) float64 {
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

	for !last {
		if t+si.step > si.end {
			last = true
			si.step = si.end - t
		}

		k1 := circ.GetDerivative()
		tmp := circ.Clone().ApplyDerivative(si.step/2, k1).GetDerivative()
		k2 := circ.Clone().ApplyDerivative(si.step, k1.WeighCopy(1.0/2).Add(1.0/2, tmp)).GetDerivative()

		kn := k1.WeighCopy(1.0/2).Add(1.0/2, k2)
		if !circ.CheckDerivative(si.step, kn) {
			si.step = circ.CalculateOptimalStep(si.step, kn)

			last = true
		}

		circ.ApplyDerivative(si.step, kn)
		t += si.step

		si.saveFn(t, circ)
	}

	span.SetTag("finish-point", t)
	span.SetTag("finish-step", si.step)

	return t
}
