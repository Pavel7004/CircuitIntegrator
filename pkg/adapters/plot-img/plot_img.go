package graph

import (
	"context"
	"image/color"

	"github.com/Pavel7004/Common/tracing"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type InfoPlotter struct {
	plot       *plot.Plot
	color      color.Color
	points     plotter.XYs
	bufferSize int
}

func NewInfoPlotter(bufferSize int) *InfoPlotter {
	p := plot.New()

	p.Add(plotter.NewGrid())
	p.X.Label.Text = "t"
	p.Y.Label.Text = "x(t)"

	return &InfoPlotter{
		plot:       p,
		color:      nil,
		points:     make(plotter.XYs, 0, bufferSize),
		bufferSize: bufferSize,
	}
}

func (ip *InfoPlotter) SaveToFile(ctx context.Context, filename string) error {
	span, _ := tracing.StartSpanFromContext(ctx)
	span.SetTag("filename", filename)
	span.SetTag("buffer-size", ip.bufferSize)

	defer span.Finish()

	if len(ip.points) > 1 {
		ip.plotPoints()
	}

	w, h := 200.0, 200.0

	c := canvas.New(w, h)
	gonumC := renderers.NewGonumPlot(c)

	ip.plot.Draw(gonumC)

	return renderers.Write(filename, c)
}

func (ip *InfoPlotter) SetPlotColor(color color.Color) {
	ip.color = color
}

func (ip *InfoPlotter) SetYLabel(label string) {
	ip.plot.Y.Label.Text = label
}

func (ip *InfoPlotter) PlotFunc(color color.Color, fn func(x float64) float64) {
	pFn := plotter.NewFunction(fn)

	pFn.Color = color
	pFn.Samples = 500

	ip.plot.Add(pFn)
}

func (ip *InfoPlotter) AddPoint(x, y float64) {
	ip.points = append(ip.points, plotter.XY{X: x, Y: y})

	if len(ip.points) == ip.bufferSize {
		ip.plotPoints()
	}
}

func (ip *InfoPlotter) plotPoints() {
	l, err := plotter.NewLine(ip.points)
	if err != nil {
		panic(err)
	}

	l.LineStyle.Color = ip.color
	ip.plot.Add(l)

	lastPoint := ip.points[len(ip.points)-1]
	ip.points = make(plotter.XYs, 0, ip.bufferSize)
	ip.points = append(ip.points, lastPoint)
}
