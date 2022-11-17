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
	"github.com/spf13/cobra"

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

	err := os.Mkdir(s.FolderName, os.ModePerm)
	if os.IsNotExist(err) {
		cobra.CheckErr(err)
	}

	for _, int := range domain.Integrators {
		p.wg.Add(3)
		go p.plotSingleTrigger(ctx, int)
		go p.plotDiffSingleTrigger(ctx, int)
		go p.plotMultiTrigger(ctx, int)
	}

	p.wg.Wait()
}

func (p *PlotterCli) plotSingleTrigger(ctx context.Context, int integrator.NewIntFunc) {
	s := p.Settings

	gr := plotter.NewInfoPlotter(s.BuffSize)

	pointgenerator.Generate(ctx, &pointgenerator.Args{
		Circuit: p.Circuit,
		Step:    s.Step,
		SaveFn: func(t float64, x *circuit.Circuit) error {
			gr.AddPoint(t, x.GetLoadVoltage())

			return nil
		},
		NewIntFn: int,
	})
	gr.PlotFunc(color.RGBA{R: 255, A: 255}, p.Circuit.GetLoadVoltageFunc())

	cobra.CheckErr(gr.SaveToFile(ctx, path.Join(s.FolderName, misc.GetFuncModule(int)+"_theory."+p.Settings.Format)))

	p.wg.Done()
}

func (p *PlotterCli) plotDiffSingleTrigger(ctx context.Context, int integrator.NewIntFunc) {
	s := p.Settings

	gr := plotter.NewInfoPlotter(s.BuffSize)
	gr.SetYLabel("X(t), %")

	theory := p.Circuit.GetLoadVoltageFunc()

	pointgenerator.Generate(ctx, &pointgenerator.Args{
		Circuit: p.Circuit,
		Step:    s.Step,
		SaveFn: func(t float64, x *circuit.Circuit) error {
			vol := x.GetLoadVoltage()
			if vol <= 0 {
				gr.AddPoint(t, 0.0)
			} else {
				gr.AddPoint(t, math.Abs(vol-theory(t))/vol*100)
			}

			return nil
		},
		NewIntFn: int,
	})

	cobra.CheckErr(gr.SaveToFile(ctx, path.Join(s.FolderName, misc.GetFuncModule(int)+"_diffErr."+p.Settings.Format)))

	p.wg.Done()
}

func (p *PlotterCli) plotMultiTrigger(ctx context.Context, int integrator.NewIntFunc) {
	s := p.Settings

	gr := plotter.NewInfoPlotter(s.BuffSize)

	ctx = context.WithValue(ctx, pointgenerator.EndPoint, 200.0)
	pointgenerator.Generate(ctx, &pointgenerator.Args{
		Circuit: p.Circuit,
		Step:    s.Step,
		SaveFn: func(t float64, x *circuit.Circuit) error {
			gr.AddPoint(t, x.GetLoadVoltage())

			return nil
		},
		NewIntFn: int,
	})

	cobra.CheckErr(gr.SaveToFile(ctx, path.Join(s.FolderName, misc.GetFuncModule(int)+"_multiTicks."+p.Settings.Format)))

	p.wg.Done()
}
