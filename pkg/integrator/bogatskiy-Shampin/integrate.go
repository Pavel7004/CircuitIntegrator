package bogatskiyshampin

import (
	"context"

	"github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/common"
	"github.com/Pavel7004/GraphPlot/pkg/integrator"
	"github.com/opentracing/opentracing-go"
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

func (si *ShapinInt) Integrate(ctx context.Context, circ *circuit.Circuit) {
	span, ctx := opentracing.StartSpanFromContext(ctx, common.GetFuncName())
	span.SetTag("StartPoint", si.begin)
	span.SetTag("EndPoint", si.end)
	span.SetTag("Step", si.step)
	span.SetTag("RK-stages", 3)

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
		k2.ApplyDerivative(si.step/2, k1.GetDerivative())
		k3 := circ.Clone()
		k3.ApplyDerivative(si.step*3/4, k2.GetDerivative())

		circ.ApplyDerivative(si.step*2/9, k1.GetDerivative())
		circ.ApplyDerivative(si.step/3, k2.GetDerivative())
		circ.ApplyDerivative(si.step*4/9, k3.GetDerivative())

		t += si.step
		si.saveFn(t, circ)
	}
}
