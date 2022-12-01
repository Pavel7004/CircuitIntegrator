package euler

import (
	"context"

	"github.com/Pavel7004/Common/tracing"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator"
)

type EulerInt struct {
	begin  float64
	end    float64
	step   float64
	saveFn func(t float64, x *circuit.Circuit) error
}

var _ integrator.Integrator = (*EulerInt)(nil)

func NewEulerInt(begin, end, step float64, saveFn func(t float64, x *circuit.Circuit) error) integrator.Integrator {
	return &EulerInt{
		begin:  begin,
		end:    end,
		step:   step,
		saveFn: saveFn,
	}
}

func (ei *EulerInt) Integrate(ctx context.Context, circ *circuit.Circuit) float64 {
	span, _ := tracing.StartSpanFromContext(ctx)
	span.SetTag("start-point", ei.begin)
	span.SetTag("end-point", ei.end)
	span.SetTag("step", ei.step)
	span.SetTag("rk-stages", 1)

	defer span.Finish()

	var (
		t = ei.begin

		last bool
	)

	for !last {
		if t+ei.step > ei.end {
			last = true
			ei.step = ei.end - t
		}

		k1 := circ.GetDerivative()

		if !circ.CheckDerivative(ei.step, k1) {
			ei.step = circ.CalculateOptimalStep(ei.step, k1)

			last = true
		}

		circ.ApplyDerivative(ei.step, k1)
		t += ei.step

		if err := ei.saveFn(t, circ); err != nil {
			span.SetTag("Error", err)
			break
		}
	}

	span.SetTag("finish-point", t)
	span.SetTag("finish-step", ei.step)

	return t
}
