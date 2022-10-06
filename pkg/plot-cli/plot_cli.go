package program

import (
	"context"
	"image/color"
	"math"
	"os"
	"path"

	misc "github.com/Pavel7004/Common/misc"
	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator"
	bogatskiyshampin "github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/bogatskiy-Shampin"
	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/euler"
	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/midpoint"
	midpointimplicit "github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/midpoint-implicit"
	threeeighth "github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/three-eighth"
	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/trapeziod"
	plotter "github.com/Pavel7004/GraphPlot/pkg/adapter/plotter"
	"github.com/Pavel7004/GraphPlot/pkg/circuit"
)

func Run(ctx context.Context, circ *circuit.Circuit, step float64, folderName string, buffSize, dpi int) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	integrators := []integrator.NewIntFunc{
		euler.NewEulerInt,
		midpoint.NewMidpointInt,
		midpointimplicit.NewMidpointImplInt,
		bogatskiyshampin.NewShampinInt,
		threeeighth.NewThreeEighthInt,
		trapeziod.NewTrapezoidInt,
	}

	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		panic(err)
	}

	for _, int := range integrators {
		gr := plotter.NewInfoPlotter(buffSize, dpi)

		plotSystem(ctx, gr, circ, step, int)
		plotTheory(ctx, gr, circ)

		gr.SaveToFile(ctx, path.Join(folderName, misc.GetFuncModule(int)+"_theory.png"))
	}

	for _, int := range integrators {
		gr := plotter.NewInfoPlotter(buffSize, dpi)

		plotDiffFunc(ctx, gr, circ, step, int)

		gr.SaveToFile(ctx, path.Join(folderName, misc.GetFuncModule(int)+"_diffErr.png"))
	}

	for _, int := range integrators {
		gr := plotter.NewInfoPlotter(buffSize, dpi)

		ctx := context.WithValue(ctx, "end", 200.0)
		plotSystem(ctx, gr, circ, step, int)

		gr.SaveToFile(ctx, path.Join(folderName, misc.GetFuncModule(int)+"_multiTicks.png"))
	}
}

func plotSystem(ctx context.Context, gr *plotter.InfoPlotter, circ *circuit.Circuit, step float64, newInt integrator.NewIntFunc) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		st    = circ.Clone()
		left  = 0.0
		right float64
	)

	right, ok := ctx.Value("end").(float64)
	if !ok {
		right = 60
	}

	for left < right {
		int := newInt(left, right, step, func(t float64, x *circuit.Circuit) {
			gr.AddPoint(t, x.GetLoadVoltage())
		})

		left = int.Integrate(ctx, st)

		st.ToggleState()
	}
}

func plotTheory(ctx context.Context, gr *plotter.InfoPlotter, circ *circuit.Circuit) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	st := circ.Clone()
	gr.PlotFunc(color.RGBA{R: 255, A: 255}, st.GetLoadVoltageFunc())
}

func plotDiffFunc(ctx context.Context, gr *plotter.InfoPlotter, circ *circuit.Circuit, step float64, newInt integrator.NewIntFunc) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		st     = circ.Clone()
		theory = st.GetLoadVoltageFunc()
		left   = 0.0
		right  = 60.0
	)

	gr.SetYLabel("x(t), %")
	for left < right {
		int := newInt(left, right, step, func(t float64, x *circuit.Circuit) {
			vol := x.GetLoadVoltage()
			if vol < 0.0001 {
				gr.AddPoint(t, 0.0)
			} else {
				gr.AddPoint(t, math.Abs(vol-theory(t))/vol*100)
			}
		})

		left = int.Integrate(ctx, st)

		st.ToggleState()
	}
}
