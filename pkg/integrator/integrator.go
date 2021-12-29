package integrator

import "github.com/Pavel7004/GraphPlot/pkg/circuit"

type Integrator interface {
	Integrate(st *circuit.Circuit)
}
