package main

import (
	"context"
	"image/color"
	"io"
	"math"

	. "github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/cli"
	"github.com/Pavel7004/GraphPlot/pkg/graph"
	"github.com/Pavel7004/GraphPlot/pkg/integrator"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"

	. "github.com/Pavel7004/GraphPlot/pkg/integrator/three-eighth"
)

func main() {
	cli.ParseArgs()
	closer := InitTracing()
	defer closer.Close()

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "main")
	span.SetTag("Step", cli.Step)
	span.SetTag("NumberOfCapacitors", cli.CapCount)
	span.SetTag("Dpi", cli.Dpi)
	span.SetTag("Filename", cli.Filename)

	defer span.Finish()

	chargeCirc := ChargeComponents{
		SupplyVoltage:     6000,
		Capacity:          0.001,
		Resistance:        5000,
		StagesCount:       cli.CapCount,
		GapTriggerVoltage: 5700,
		HoldingVoltage:    1,
	}

	load := LoadComponents{
		Resistance: cli.LoadRes,
	}

	circ := NewCircuit(chargeCirc, load)

	gr := graph.NewInfoPlotter(cli.PointBuffSize, cli.Dpi)
	// PlotTheory(ctx, gr, circ)
	PlotSystem(ctx, gr, circ, NewThreeEighthInt)
	gr.SaveToFile(ctx, cli.Filename)
}

func PlotSystem(ctx context.Context, gr *graph.InfoPlotter, circ *Circuit, newInt integrator.NewIntFunc) {
	var (
		st     = circ.Clone()
		period = st.GetSystemPeriod()
		left   = 0.0
		right  = period
	)

	for right <= 200 {
		int := newInt(left, right, cli.Step, func(t float64, x *Circuit) {
			gr.AddPoint(t, x.GetLoadVoltage())
		})

		int.Integrate(ctx, st)
		st.ToggleState()
		left = right
		right += period
	}
}

func PlotTheory(ctx context.Context, gr *graph.InfoPlotter, circ *Circuit) {
	st := circ.Clone()
	gr.PlotFunc(color.RGBA{R: 255, A: 255}, st.GetLoadVoltageFunc())
}

func PlotDiffFunc(ctx context.Context, gr *graph.InfoPlotter, circ *Circuit, newInt integrator.NewIntFunc) {
	var (
		st     = circ.Clone()
		period = st.GetSystemPeriod()
		left   = 0.0
		right  = period
		theory = st.GetLoadVoltageFunc()
	)

	gr.SetYLabel("x(t), %")
	for right <= 60 {
		int := newInt(left, right, cli.Step, func(t float64, x *Circuit) {
			vol := x.GetLoadVoltage()
			if vol < 0.0001 {
				gr.AddPoint(t, 0.0)
			} else {
				gr.AddPoint(t, math.Abs(vol-theory(t))/vol*100)
			}
		})

		int.Integrate(ctx, st)
		st.ToggleState()
		left = right
		right += period
	}
}

func InitTracing() io.Closer {
	cfg := jaegercfg.Configuration{
		ServiceName: "GraphPlot",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	jMetricsFactory := metrics.NullFactory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(nil),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	return closer
}
