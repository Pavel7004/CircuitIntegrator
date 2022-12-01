package trapeziod_test

import (
	"context"
	"testing"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator/trapeziod"
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

func BenchmarkTrapezoid(b *testing.B) {
	int := trapeziod.NewTrapezoidInt(0, 60, 0.0001, func(t float64, x *circuit.Circuit) error {
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
