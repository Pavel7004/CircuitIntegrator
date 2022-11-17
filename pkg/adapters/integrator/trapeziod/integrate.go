package trapeziod

import (
	"context"

	"github.com/Pavel7004/Common/tracing"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator"
)

const errorTolerance = 1e-4

type TrapezoidInt struct {
	begin  float64
	end    float64
	step   float64
	saveFn func(t float64, x *circuit.Circuit) error
}

var _ integrator.Integrator = (*TrapezoidInt)(nil)

func NewTrapezoidInt(begin, end, step float64, saveFn func(t float64, x *circuit.Circuit) error) integrator.Integrator {
	return &TrapezoidInt{
		begin:  begin,
		end:    end,
		step:   step,
		saveFn: saveFn,
	}
}

func (int *TrapezoidInt) Integrate(ctx context.Context, circ *circuit.Circuit) float64 {
	span, _ := tracing.StartSpanFromContext(ctx)
	span.SetTag("StartPoint", int.begin)
	span.SetTag("EndPoint", int.end)
	span.SetTag("Step", int.step)
	span.SetTag("RK-stages", 2)

	defer span.Finish()

	var (
		t = int.begin

		last bool
	)

	for !last {
		if t+int.step > int.end {
			last = true
			int.step = int.end - t
		}

		k1 := circ.GetDerivative()
		tmp := circ.Clone().ApplyDerivative(int.step/2, k1).GetDerivative()
		k2 := circ.Clone().ApplyDerivative(int.step, k1.WeighCopy(1.0/2).Add(1.0/2, tmp)).GetDerivative()

		kn := k1.WeighCopy(1.0/2).Add(1.0/2, k2)
		if !circ.CheckDerivative(int.step, kn) {
			int.step = circ.CalculateOptimalStep(int.step, kn)

			last = true
		}

		err := 100.0
		prev := circ.Clone()
		for i := 0; i < 10 && err > errorTolerance; i++ {
			k1 := circ.GetDerivative()
			tmp := circ.Clone().ApplyDerivative(int.step/2, k1).GetDerivative()
			k2 := circ.Clone().ApplyDerivative(int.step, k1.WeighCopy(1.0/2).Add(1.0/2, tmp)).GetDerivative()

			kn := k1.WeighCopy(1.0/2).Add(1.0/2, k2)
			err = circ.ImplicitStep(int.step, kn, prev)
		}

		t += int.step

		if err := int.saveFn(t, circ); err != nil {
			span.SetTag("Error", err)
			break
		}
	}

	span.SetTag("finish-point", t)
	span.SetTag("finish-step", int.step)

	return t
}
