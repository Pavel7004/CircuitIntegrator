package main

import (
	"context"
	"io"

	"github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/GraphPlot/pkg/adapter/cli"
	. "github.com/Pavel7004/GraphPlot/pkg/circuit"
	plotcli "github.com/Pavel7004/GraphPlot/pkg/plot-cli"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

func main() {
	closer := InitTracing()
	defer closer.Close()

	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	cli.ParseArgs()

	span.SetTag("Step", cli.Step)
	span.SetTag("NumberOfCapacitors", cli.CapCount)
	span.SetTag("Dpi", cli.Dpi)
	span.SetTag("Dirname", cli.DirName)

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
	p := plotcli.NewPlotterCli(circ, &plotcli.Settings{
		Step:       cli.Step,
		FolderName: cli.DirName,
		BuffSize:   cli.PointBuffSize,
		Dpi:        cli.Dpi,
	})

	p.Plot()
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
