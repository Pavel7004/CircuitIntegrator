package bogatskiyshampin

import (
	"context"

	"github.com/Pavel7004/Common/tracing"

	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator"
	"github.com/Pavel7004/GraphPlot/pkg/circuit"
)

type ShapinInt struct {
	begin  float64
	end    float64
	step   float64
	saveFn func(t float64, x *circuit.Circuit)
}

var _ integrator.Integrator = (*ShapinInt)(nil)

func NewShampinInt(begin, end, step float64, saveFn func(t float64, x *circuit.Circuit)) integrator.Integrator {
	return &ShapinInt{
		begin:  begin,
		end:    end,
		step:   step,
		saveFn: saveFn,
	}
}

func (si *ShapinInt) Integrate(ctx context.Context, circ *circuit.Circuit) float64 {
	span, _ := tracing.StartSpanFromContext(ctx)
	span.SetTag("StartPoint", si.begin)
	span.SetTag("EndPoint", si.end)
	span.SetTag("Step", si.step)
	span.SetTag("RK-stages", 3)

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
		k2 := circ.Clone().ApplyDerivative(si.step/2, k1).GetDerivative()
		k3 := circ.Clone().ApplyDerivative(3*si.step/4, k2).GetDerivative()

		kn := k1.WeighCopy(2.0/9).Add(1.0/3, k2).Add(4.0/9, k3)
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
