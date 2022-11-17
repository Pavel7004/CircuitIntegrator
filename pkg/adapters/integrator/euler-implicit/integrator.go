package eulerimplicit

import (
	"context"

	"github.com/Pavel7004/Common/tracing"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator"
)

const errorTolerance = 1e-4

type EulerImplInt struct {
	begin  float64
	end    float64
	step   float64
	saveFn func(t float64, x *circuit.Circuit) error
}

var _ integrator.Integrator = (*EulerImplInt)(nil)

func NewEulerImplInt(begin, end, step float64, saveFn func(t float64, x *circuit.Circuit) error) integrator.Integrator {
	return &EulerImplInt{
		begin:  begin,
		end:    end,
		step:   step,
		saveFn: saveFn,
	}
}

func (int *EulerImplInt) Integrate(ctx context.Context, circ *circuit.Circuit) float64 {
	span, _ := tracing.StartSpanFromContext(ctx)
	span.SetTag("start-point", int.begin)
	span.SetTag("end-point", int.end)
	span.SetTag("step", int.step)

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

		if !circ.CheckDerivative(int.step, k1) {
			int.step = circ.CalculateOptimalStep(int.step, k1)

			last = true
		}

		err := 100.0
		for i := 0; i < 10 && err > errorTolerance; i++ {
			prev := circ.Clone()
			d := circ.Clone().ApplyDerivative(int.step, circ.GetDerivative()).GetDerivative()
			err = circ.ImplicitStep(int.step, d, prev)
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
