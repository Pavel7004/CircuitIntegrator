package main

import (
	"context"
	"io"

	. "github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/cli"
	"github.com/Pavel7004/GraphPlot/pkg/common/tracing"
	"github.com/Pavel7004/GraphPlot/pkg/program"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

func main() {
	closer := InitTracing()
	defer closer.Close()

	span, ctx := tracing.StartSpanFromContext(context.Background())
	span.SetTag("Step", cli.Step)
	span.SetTag("NumberOfCapacitors", cli.CapCount)
	span.SetTag("Dpi", cli.Dpi)
	span.SetTag("Dirname", cli.DirName)

	defer span.Finish()

	cli.ParseArgs()

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

	program.Run(ctx, circ, cli.DirName, cli.PointBuffSize, cli.Dpi)
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
