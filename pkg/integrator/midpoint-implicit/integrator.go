package midpointimplicit

import (
	"context"

	"github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/common/tracing"
	"github.com/Pavel7004/GraphPlot/pkg/integrator"
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

func (si *MidpointImpInt) Integrate(ctx context.Context, circ *circuit.Circuit) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	span.SetTag("StartPoint", si.begin)
	span.SetTag("EndPoint", si.end)
	span.SetTag("Step", si.step)
	span.SetTag("RK-stages", 2)

	defer span.Finish()

	var (
		t         = si.begin
		prevDeriv = si.GetDerivativeWithCoeff(circ.GetDerivative(), 1.0/2)
		last      bool
	)

	for !last {
		if t+si.step > si.end {
			last = true
			si.step = si.end - t
		}

		currDeriv := si.GetDerivativeWithCoeff(circ.GetDerivative(), 1.0/2)
		k1 := si.GetSumOfDerivatives(currDeriv, prevDeriv)
		prevDeriv = currDeriv

		circ.ApplyDerivative(si.step, k1)

		t += si.step
		si.saveFn(t, circ)
	}
}

func (si *MidpointImpInt) GetDerivativeWithCoeff(der []float64, k float64) []float64 {
	res := make([]float64, len(der))
	for i := range der {
		res[i] = der[i] * k
	}
	return res
}

func (si *MidpointImpInt) GetSumOfDerivatives(der1, der2 []float64) []float64 {
	res := make([]float64, len(der1))
	for i := range res {
		res[i] = der1[i] + der2[i]
	}
	return res
}
