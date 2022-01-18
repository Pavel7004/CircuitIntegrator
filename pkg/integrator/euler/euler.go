package euler

import (
	"github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/integrator"
)

type EulerInt struct {
	begin  float64
	end    float64
	step   float64
	saveFn func(t float64, x *circuit.Circuit)
}

var _ integrator.Integrator = (*EulerInt)(nil)

func NewEulerInt(begin, end, step float64, saveFn func(t float64, x *circuit.Circuit)) integrator.Integrator {
	return &EulerInt{
		begin:  begin,
		end:    end,
		step:   step,
		saveFn: saveFn,
	}
}

func (ei *EulerInt) Integrate(circ *circuit.Circuit) {
	var (
		t    = ei.begin
		last bool
	)
	for !last {
		if t+ei.step > ei.end {
			last = true
			ei.step = ei.end - t
		}
		circ.ApplyDerivative(ei.step, circ.GetDerivative())
		t += ei.step
		ei.saveFn(t, circ)
	}
}
