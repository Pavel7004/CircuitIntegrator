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

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/adapters/integrator"
	plotter "github.com/Pavel7004/GraphPlot/pkg/adapters/plot-img"
	pointgenerator "github.com/Pavel7004/GraphPlot/pkg/components/point-generator"
	"github.com/Pavel7004/GraphPlot/pkg/domain"
)

type PlotterCli struct {
	Settings *Settings
	Circuit  *circuit.Circuit

	wg *sync.WaitGroup
}

func NewPlotterCli(circuit *circuit.Circuit, settings *Settings) *PlotterCli {
	p := new(PlotterCli)

	p.Circuit = circuit
	p.Settings = settings
	p.wg = new(sync.WaitGroup)

	return p
}

func (p *PlotterCli) Plot() {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	s := p.Settings

	if err := os.MkdirAll(s.FolderName, os.ModePerm); err != nil {
		panic(err)
	}

	for _, int := range domain.Integrators {
		p.wg.Add(3)
		go p.PlotSingleTrigger(ctx, int)
		go p.PlotDiffSingleTrigger(ctx, int)
		go p.PlotMultiTrigger(ctx, int)
	}

	p.wg.Wait()
}

func (p *PlotterCli) PlotSingleTrigger(ctx context.Context, int integrator.NewIntFunc) {
	s := p.Settings

	gr := plotter.NewInfoPlotter(s.BuffSize, s.Dpi)

	pointgenerator.GeneratePoints(ctx, &pointgenerator.Args{
		Circuit: p.Circuit,
		Step:    s.Step,
		SaveFn: func(t float64, x *circuit.Circuit) error {
			gr.AddPoint(t, x.GetLoadVoltage())

			return nil
		},
		NewIntFn: int,
	})
	gr.PlotFunc(color.RGBA{R: 255, A: 255}, p.Circuit.GetLoadVoltageFunc())

	gr.SaveToFile(ctx, path.Join(s.FolderName, misc.GetFuncModule(int)+"_theory.svg"))

	p.wg.Done()
}

func (p *PlotterCli) PlotDiffSingleTrigger(ctx context.Context, int integrator.NewIntFunc) {
	s := p.Settings

	gr := plotter.NewInfoPlotter(s.BuffSize, s.Dpi)
	gr.SetYLabel("X(t), %")

	theory := p.Circuit.GetLoadVoltageFunc()

	pointgenerator.GeneratePoints(ctx, &pointgenerator.Args{
		Circuit: p.Circuit,
		Step:    s.Step,
		SaveFn: func(t float64, x *circuit.Circuit) error {
			vol := x.GetLoadVoltage()
			if vol < 0.0001 {
				gr.AddPoint(t, 0.0)
			} else {
				gr.AddPoint(t, math.Abs(vol-theory(t))/vol*100)
			}

			return nil
		},
		NewIntFn: int,
	})

	gr.SaveToFile(ctx, path.Join(s.FolderName, misc.GetFuncModule(int)+"_diffErr.svg"))

	p.wg.Done()
}

func (p *PlotterCli) PlotMultiTrigger(ctx context.Context, int integrator.NewIntFunc) {
	s := p.Settings

	gr := plotter.NewInfoPlotter(s.BuffSize, s.Dpi)

	ctx = context.WithValue(ctx, pointgenerator.EndPoint, 200.0)
	pointgenerator.GeneratePoints(ctx, &pointgenerator.Args{
		Circuit: p.Circuit,
		Step:    s.Step,
		SaveFn: func(t float64, x *circuit.Circuit) error {
			gr.AddPoint(t, x.GetLoadVoltage())

			return nil
		},
		NewIntFn: int,
	})

	gr.SaveToFile(ctx, path.Join(s.FolderName, misc.GetFuncModule(int)+"_multiTicks.svg"))

	p.wg.Done()
}
