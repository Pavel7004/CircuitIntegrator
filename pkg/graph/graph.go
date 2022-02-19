package graph

import (
	"context"
	"image"
	"image/color"
	"os"

	"github.com/Pavel7004/GraphPlot/pkg/common/tracing"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"gonum.org/v1/plot/vg/vgsvg"
)

type InfoPlotter struct {
	plot       *plot.Plot
	dpi        int
	color      color.Color
	points     plotter.XYs
	bufferSize int
}

func NewInfoPlotter(bufferSize, dpi int) *InfoPlotter {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Add(plotter.NewGrid())
	p.X.Label.Text = "t"
	p.Y.Label.Text = "x(t)"

	return &InfoPlotter{
		plot:       p,
		dpi:        dpi,
		color:      nil,
		points:     make(plotter.XYs, 0, bufferSize),
		bufferSize: bufferSize,
	}
}

func (ip *InfoPlotter) EnableLogScale() {
	ip.plot.Y.Scale = plot.LogScale{}
}

func (ip *InfoPlotter) DrawInImage(ctx context.Context) image.Image {
	span, ctx := tracing.StartSpanFromContext(ctx)
	span.SetTag("buffer size", ip.bufferSize)

	defer span.Finish()

	if len(ip.points) != 0 {
		ip.plotPoints()
	}

	img := image.NewRGBA(image.Rect(0, 0, 16*ip.dpi, 16*ip.dpi))
	c := vgimg.NewWith(vgimg.UseImage(img))
	ip.plot.Draw(draw.New(c))

	return c.Image()
}

func (ip *InfoPlotter) WriteSVGToStdout(ctx context.Context) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	span.SetTag("buffer size", ip.bufferSize)

	defer span.Finish()

	if len(ip.points) != 0 {
		ip.plotPoints()
	}

	c := vgsvg.New(3*vg.Inch, 3*vg.Inch)

	ip.plot.Draw(draw.New(c))
	if _, err := c.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
}

func (ip *InfoPlotter) SaveToFile(ctx context.Context, filename string) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	span.SetTag("filename", filename)
	span.SetTag("buffer size", ip.bufferSize)

	defer span.Finish()

	if len(ip.points) != 0 {
		ip.plotPoints()
	}

	if err := ip.plot.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		panic(err)
	}
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
	pFn.Width = vg.Points(1)
	pFn.Samples = 200

	ip.plot.Add(pFn)
}

func (ip *InfoPlotter) plotPoints() {
	l, err := plotter.NewLine(ip.points)
	if err != nil {
		panic(err)
	}

	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Color = ip.color
	ip.plot.Add(l)

	lastPoint := ip.points[len(ip.points)-1]
	ip.points = make(plotter.XYs, 0, ip.bufferSize)
	ip.points = append(ip.points, lastPoint)
}

func (ip *InfoPlotter) AddPoint(x, y float64) {
	ip.points = append(ip.points, plotter.XY{X: x, Y: y})

	if len(ip.points) == ip.bufferSize {
		ip.plotPoints()
	}
}
