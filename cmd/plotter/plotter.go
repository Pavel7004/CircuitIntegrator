package main

import (
	"context"

	common "github.com/Pavel7004/Common/tracing"
	"github.com/Pavel7004/GraphPlot/pkg/adapter/cli"
	. "github.com/Pavel7004/GraphPlot/pkg/components/circuit"
	plotcli "github.com/Pavel7004/GraphPlot/pkg/components/plot-cli"
	"github.com/Pavel7004/GraphPlot/pkg/infra/tracing"
)

func main() {
	closer := tracing.Init()
	defer closer.Close()

	span, ctx := common.StartSpanFromContext(context.Background())
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

	p.Plot(ctx)
}
