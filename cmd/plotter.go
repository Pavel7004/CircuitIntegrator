package main

import (
	"image/color"
	"image/png"
	"os"

	"github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/graph"
	. "github.com/Pavel7004/GraphPlot/pkg/integrator/euler"
	. "github.com/Pavel7004/GraphPlot/pkg/integrator/midpoint"
)

func main() {
	chargeCirc := circuit.ChargeComponents{
		SupplyVoltage:     5000,
		Capacity:          0.001,
		Resistance:        5000,
		StagesCount:       5,
		GapTriggerVoltage: 3800,
	}
	load := circuit.LoadComponents{
		Resistance: 10000,
	}
	st := circuit.NewCircuit(chargeCirc, load)
	gr := graph.NewInfoPlotter(40)
	int := NewEulerInt(0, 60, 0.1, func(t float64, x *circuit.Circuit) {
		gr.AddPoint(t, x.GetLoadVoltage())
	})
	intMid := NewMidpointInt(0, 60, 0.1, func(t float64, x *circuit.Circuit) {
		gr.AddPoint(t, x.GetLoadVoltage())
	})
	int.Integrate(st)
	gr.PrepareToAddNewPlot(color.RGBA{B: 255, A: 255})
	st = circuit.NewCircuit(chargeCirc, load)
	intMid.Integrate(st)
	gr.PlotFunc(color.RGBA{R: 255, A: 255}, st.GetLoadVoltageFunc())
	imgFile, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}
	if err := png.Encode(imgFile, gr.DrawInImage()); err != nil {
		panic(err)
	}
	imgFile.Close()
}