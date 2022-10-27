package plotcli

import (
	"context"
	"image/color"
	"math"
	"os"
	"path"
	"sync"

	misc "github.com/Pavel7004/Common/misc"
	"github.com/Pavel7004/Common/tracing"

	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator"
	bogatskiyshampin "github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/bogatskiy-Shampin"
	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/euler"
	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/midpoint"
	midpointimplicit "github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/midpoint-implicit"
	threeeighth "github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/three-eighth"
	"github.com/Pavel7004/GraphPlot/pkg/adapter/integrator/trapeziod"
	plotter "github.com/Pavel7004/GraphPlot/pkg/adapter/plot-img"
	"github.com/Pavel7004/GraphPlot/pkg/circuit"
	pointgenerator "github.com/Pavel7004/GraphPlot/pkg/point-generator"
)

type PlotterCli struct {
	Settings *Settings
	Circuit  *circuit.Circuit

	wg          *sync.WaitGroup
	integrators []integrator.NewIntFunc
}

func NewPlotterCli(circuit *circuit.Circuit, settings *Settings) *PlotterCli {
	p := new(PlotterCli)

	p.Circuit = circuit
	p.Settings = settings
	p.wg = new(sync.WaitGroup)

	p.integrators = []integrator.NewIntFunc{
		euler.NewEulerInt,
		midpoint.NewMidpointInt,
		midpointimplicit.NewMidpointImplInt,
		bogatskiyshampin.NewShampinInt,
		threeeighth.NewThreeEighthInt,
		trapeziod.NewTrapezoidInt,
	}

	return p
}

func (p *PlotterCli) Plot(ctx context.Context) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	s := p.Settings

	if err := os.MkdirAll(s.FolderName, os.ModePerm); err != nil {
		panic(err)
	}

	for _, int := range p.integrators {
		p.wg.Add(3)
		p.PlotSingleTrigger(ctx, int)
		p.PlotDiffSingleTrigger(ctx, int)
		p.PlotMultiTrigger(ctx, int)
	}

	p.wg.Wait()
}

func (p *PlotterCli) PlotSingleTrigger(ctx context.Context, int integrator.NewIntFunc) {
	s := p.Settings

	gr := plotter.NewInfoPlotter(s.BuffSize, s.Dpi)

	pointgenerator.GeneratePoints(ctx, &pointgenerator.Args{
		Circuit: p.Circuit,
		Step:    s.Step,
		SaveFn: func(t float64, x *circuit.Circuit) {
			gr.AddPoint(t, x.GetLoadVoltage())
		},
		NewIntFn: int,
	})
	gr.PlotFunc(color.RGBA{R: 255, A: 255}, p.Circuit.GetLoadVoltageFunc())

	gr.SaveToFile(ctx, path.Join(s.FolderName, misc.GetFuncModule(int)+"_theory.png"))

	p.wg.Done()
}

func (p *PlotterCli) PlotDiffSingleTrigger(ctx context.Context, int integrator.NewIntFunc) {
	s := p.Settings

	gr := plotter.NewInfoPlotter(s.BuffSize, s.Dpi)
	gr.SetYLabel("X(t) %")

	theory := p.Circuit.GetLoadVoltageFunc()

	pointgenerator.GeneratePoints(ctx, &pointgenerator.Args{
		Circuit: p.Circuit,
		Step:    s.Step,
		SaveFn: func(t float64, x *circuit.Circuit) {
			vol := x.GetLoadVoltage()
			if vol < 0.0001 {
				gr.AddPoint(t, 0.0)
			} else {
				gr.AddPoint(t, math.Abs(vol-theory(t))/vol*100)
			}
		},
		NewIntFn: int,
	})

	gr.SaveToFile(ctx, path.Join(s.FolderName, misc.GetFuncModule(int)+"_diffErr.png"))

	p.wg.Done()
}

func (p *PlotterCli) PlotMultiTrigger(ctx context.Context, int integrator.NewIntFunc) {
	s := p.Settings

	gr := plotter.NewInfoPlotter(s.BuffSize, s.Dpi)

	ctx = context.WithValue(ctx, pointgenerator.EndPoint, 200.0)
	pointgenerator.GeneratePoints(ctx, &pointgenerator.Args{
		Circuit: p.Circuit,
		Step:    s.Step,
		SaveFn: func(t float64, x *circuit.Circuit) {
			gr.AddPoint(t, x.GetLoadVoltage())
		},
		NewIntFn: int,
	})

	gr.SaveToFile(ctx, path.Join(s.FolderName, misc.GetFuncModule(int)+"_multiTicks.png"))

	p.wg.Done()
}
