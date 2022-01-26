package main

import (
	"image/color"
	"math"

	. "github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/cli"
	"github.com/Pavel7004/GraphPlot/pkg/graph"
	"github.com/Pavel7004/GraphPlot/pkg/integrator"

	//. "github.com/Pavel7004/GraphPlot/pkg/integrator/bogatskiy-Shampin"
	// . "github.com/Pavel7004/GraphPlot/pkg/integrator/euler"
	// . "github.com/Pavel7004/GraphPlot/pkg/integrator/three-eighth"
	// . "github.com/Pavel7004/GraphPlot/pkg/integrator/midpoint"
	. "github.com/Pavel7004/GraphPlot/pkg/integrator/midpoint-implicit"
)

type NewIntFunc func(begin, end float64, step float64, saveFn func(t float64, x *Circuit)) integrator.Integrator

func main() {
	cli.ParseArgs()
	chargeCirc := &ChargeComponents{
		SupplyVoltage:     6000,
		Capacity:          0.001,
		Resistance:        5000,
		StagesCount:       cli.CapCount,
		GapTriggerVoltage: 5700,
	}
	load := &LoadComponents{
		Resistance: cli.LoadRes,
	}
	gr := graph.NewInfoPlotter(cli.Dpi)
	// PlotTheory(gr, chargeCirc, load)
	// gr.PrepareToAddNewPlot(color.RGBA{G: 255, A: 255})
	// PlotSystem(gr, chargeCirc, load, NewMidpointImplInt)
	// PlotSystem(gr, chargeCirc, load, NewEulerInt)
	// gr.PrepareToAddNewPlot(color.RGBA{A: 255})
	// PlotDiffFunc(gr, chargeCirc, load, NewThreeEighthInt)
	// gr.PrepareToAddNewPlot(color.RGBA{G: 255, A: 255})
	PlotDiffFunc(gr, chargeCirc, load, NewMidpointImplInt)
	gr.SaveToFile(cli.Filename)
}

func PlotSystem(gr *graph.InfoPlotter, chargeCirc *ChargeComponents, load *LoadComponents, newInt NewIntFunc) {
	var (
		st     = NewCircuit(*chargeCirc, *load)
		period = st.GetSystemPeriod()
		left   = 0.0
		right  = period
	)
	for right <= 60 {
		int := newInt(left, right, cli.Step, func(t float64, x *Circuit) {
			gr.AddPoint(t, x.GetLoadVoltage())
		})
		int.Integrate(st)
		st.ToggleStateMaybe()
		left = right + cli.Step
		right += period
	}
}

func PlotTheory(gr *graph.InfoPlotter, chargeCirc *ChargeComponents, load *LoadComponents) {
	st := NewCircuit(*chargeCirc, *load)
	gr.PlotFunc(color.RGBA{R: 255, A: 255}, st.GetLoadVoltageFunc())
}

func PlotDiffFunc(gr *graph.InfoPlotter, chargeCirc *ChargeComponents, load *LoadComponents, newInt NewIntFunc) {
	var (
		st     = NewCircuit(*chargeCirc, *load)
		period = st.GetSystemPeriod()
		left   = 0.0
		right  = period
		theory = st.GetLoadVoltageFunc()
	)
	for right <= 60 {
		int := newInt(left, right, cli.Step, func(t float64, x *Circuit) {
			gr.AddPoint(t, math.Abs(x.GetLoadVoltage()-theory(t)))
		})
		int.Integrate(st)
		st.ToggleStateMaybe()
		left = right + cli.Step
		right += period
	}
}
