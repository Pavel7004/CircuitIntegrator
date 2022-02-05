package threeeighth

import (
	"context"

	"github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/integrator"
	"github.com/opentracing/opentracing-go"
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

func (si *ThreeEighthInt) Integrate(ctx context.Context, circ *circuit.Circuit) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ThreeEightsInt.Integrate")
	span.SetTag("StartPoint", si.begin)
	span.SetTag("EndPoint", si.end)
	span.SetTag("Step", si.step)
	span.SetTag("RK-stages", 4)
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
		k1 := circ.Clone()
		k2 := circ.Clone()
		k3 := circ.Clone()
		k4 := circ.Clone()
		k2.ApplyDerivative(si.step/3, k1.GetDerivative())
		k3.ApplyDerivative(-si.step/3, k1.GetDerivative())
		k3.ApplyDerivative(si.step, k2.GetDerivative())
		k4.ApplyDerivative(si.step, k1.GetDerivative())
		k4.ApplyDerivative(-si.step, k2.GetDerivative())
		k4.ApplyDerivative(si.step, k3.GetDerivative())
		circ.ApplyDerivative(si.step/8, k1.GetDerivative())
		circ.ApplyDerivative(3*si.step/8, k2.GetDerivative())
		circ.ApplyDerivative(3*si.step/8, k3.GetDerivative())
		circ.ApplyDerivative(si.step/8, k4.GetDerivative())
		t += si.step
		si.saveFn(t, circ)
	}
}
