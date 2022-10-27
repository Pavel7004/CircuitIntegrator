package threeeighth

import (
	"context"

	"github.com/Pavel7004/Common/tracing"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator"
)

type ThreeEighthInt struct {
	begin  float64
	end    float64
	step   float64
	saveFn func(t float64, x *circuit.Circuit)
}

var _ integrator.Integrator = (*ThreeEighthInt)(nil)

func NewThreeEighthInt(begin, end, step float64, saveFn func(t float64, x *circuit.Circuit)) integrator.Integrator {
	return &ThreeEighthInt{
		begin:  begin,
		end:    end,
		step:   step,
		saveFn: saveFn,
	}
}

func (si *ThreeEighthInt) Integrate(ctx context.Context, circ *circuit.Circuit) float64 {
	span, _ := tracing.StartSpanFromContext(ctx)
	span.SetTag("start-point", si.begin)
	span.SetTag("end-point", si.end)
	span.SetTag("step", si.step)
	span.SetTag("rk-stages", 4)

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
		k2 := circ.Clone().ApplyDerivative(si.step/3, k1).GetDerivative()
		k3 := circ.Clone().ApplyDerivative(-si.step/3, k1).ApplyDerivative(si.step, k2).GetDerivative()
		k4 := circ.Clone().ApplyDerivative(si.step, k1).ApplyDerivative(-si.step, k2).ApplyDerivative(si.step, k3).GetDerivative()

		kn := k1.WeighCopy(1.0/8.0).Add(3.0/8.0, k2).Add(3.0/8.0, k3).Add(1.0/8.0, k4)
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
