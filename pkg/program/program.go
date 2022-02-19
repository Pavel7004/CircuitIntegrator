package program

import (
	"context"
	"image/color"
	"math"
	"os"
	"path"

	"github.com/Pavel7004/GraphPlot/pkg/circuit"
	misc "github.com/Pavel7004/GraphPlot/pkg/common/misc"
	"github.com/Pavel7004/GraphPlot/pkg/common/tracing"
	"github.com/Pavel7004/GraphPlot/pkg/graph"
	"github.com/Pavel7004/GraphPlot/pkg/integrator"
	bogatskiyshampin "github.com/Pavel7004/GraphPlot/pkg/integrator/bogatskiy-Shampin"
	"github.com/Pavel7004/GraphPlot/pkg/integrator/euler"
	"github.com/Pavel7004/GraphPlot/pkg/integrator/midpoint"
	midpointimplicit "github.com/Pavel7004/GraphPlot/pkg/integrator/midpoint-implicit"
	threeeighth "github.com/Pavel7004/GraphPlot/pkg/integrator/three-eighth"
	"github.com/Pavel7004/GraphPlot/pkg/integrator/trapeziod"
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
		gr := graph.NewInfoPlotter(buffSize, dpi)

		PlotSystem(ctx, gr, circ, step, int)
		PlotTheory(ctx, gr, circ)

		gr.SaveToFile(ctx, path.Join(folderName, misc.GetFuncModule(int)+"_theory.png"))
	}

	for _, int := range integrators {
		gr := graph.NewInfoPlotter(buffSize, dpi)

		PlotDiffFunc(ctx, gr, circ, step, int)

		gr.SaveToFile(ctx, path.Join(folderName, misc.GetFuncModule(int)+"_diffErr.png"))
	}

	for _, int := range integrators {
		gr := graph.NewInfoPlotter(buffSize, dpi)

		ctx := context.WithValue(ctx, "end", 200.0)
		PlotSystem(ctx, gr, circ, step, int)

		gr.SaveToFile(ctx, path.Join(folderName, misc.GetFuncModule(int)+"_multiTicks.png"))
	}
}

func PlotSystem(ctx context.Context, gr *graph.InfoPlotter, circ *circuit.Circuit, step float64, newInt integrator.NewIntFunc) {
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

	for left <= right {
		int := newInt(left, right, step, func(t float64, x *circuit.Circuit) {
			gr.AddPoint(t, x.GetLoadVoltage())
		})

		left = int.Integrate(ctx, st)

		st.ToggleState()
	}
}

func PlotTheory(ctx context.Context, gr *graph.InfoPlotter, circ *circuit.Circuit) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	st := circ.Clone()
	gr.PlotFunc(color.RGBA{R: 255, A: 255}, st.GetLoadVoltageFunc())
}

func PlotDiffFunc(ctx context.Context, gr *graph.InfoPlotter, circ *circuit.Circuit, step float64, newInt integrator.NewIntFunc) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		st     = circ.Clone()
		theory = st.GetLoadVoltageFunc()
		left   = 0.0
		right  float64
	)

	gr.SetYLabel("x(t), %")
	for left <= right {
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
