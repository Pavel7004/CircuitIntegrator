package pointgenerator

import (
	"context"

	"github.com/Pavel7004/Common/tracing"
)

func GeneratePoints(ctx context.Context, args *Args) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		left  float64
		right float64
	)

	right, ok := ctx.Value("end").(float64)
	if !ok {
		right = 60
	}

	for left < right {
		int := args.NewIntFn(left, right, args.Step, args.SaveFn)

		left = int.Integrate(ctx, args.Circuit)

		args.Circuit.ToggleState()
	}
}

// func PlotDiffFunc(ctx context.Context, args Args) {
// 	span, ctx := tracing.StartSpanFromContext(ctx)
// 	defer span.Finish()

// 	var (
// 		theory = args.circ.GetLoadVoltageFunc()
// 		left   = 0.0
// 		right  = 60.0
// 	)

// 	gr.SetYLabel("x(t), %")
// 	for left < right {
// 		int := args.newInt(left, right, args.step, args.saveFn)
// 		// vol := x.GetLoadVoltage()
// 		// if vol < 0.0001 {
// 		// 	gr.AddPoint(t, 0.0)
// 		// } else {
// 		// 	gr.AddPoint(t, math.Abs(vol-theory(t))/vol*100)
// 		// }

// 		left = int.Integrate(ctx, st)

// 		st.ToggleState()
// 	}
// }
