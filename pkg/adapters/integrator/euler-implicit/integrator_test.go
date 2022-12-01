package eulerimplicit_test

import (
	"context"
	"testing"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
	eulerimplicit "github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/euler-implicit"
)

var (
	chargeCirc = &circuit.ChargeComponents{
		SupplyVoltage:     6000,
		Capacity:          0.001,
		Resistance:        5000,
		StagesCount:       6,
		GapTriggerVoltage: 5700,
		HoldingVoltage:    1,
	}

	load = &circuit.LoadComponents{
		Resistance: 10000,
	}
)

func BenchmarkEulerImplicit(b *testing.B) {
	int := eulerimplicit.NewEulerImplInt(0, 60, 0.0001, func(t float64, x *circuit.Circuit) error {
		return nil
	})
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		circ := circuit.New(chargeCirc, load)
		int.Integrate(ctx, circ)
	}
}
