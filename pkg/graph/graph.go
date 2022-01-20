package graph

import (
	"image"
	"image/color"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"gonum.org/v1/plot/vg/vgsvg"
)

type InfoPlotter struct {
	plot      *plot.Plot
	dpi       int
	lastPoint *plotter.XY
	color     color.Color
}

func NewInfoPlotter(dpi int) *InfoPlotter {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	return &InfoPlotter{
		plot:      p,
		dpi:       dpi,
		lastPoint: nil,
		color:     nil,
	}
}

func (ip *InfoPlotter) DrawInImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 16*ip.dpi, 16*ip.dpi))
	c := vgimg.NewWith(vgimg.UseImage(img))
	ip.plot.Draw(draw.New(c))
	return c.Image()
}

func (ip *InfoPlotter) WriteSVGToStdout() {
	c := vgsvg.New(3*vg.Inch, 3*vg.Inch)
	ip.plot.Draw(draw.New(c))
	if _, err := c.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
}

func (ip *InfoPlotter) SaveToFile(filename string) {
	if err := ip.plot.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		panic(err)
	}
}

func (ip *InfoPlotter) PrepareToAddNewPlot(color color.Color) {
	ip.lastPoint = nil
	ip.color = color
}

func (ip *InfoPlotter) PlotFunc(color color.Color, fn func(x float64) float64) {
	pFn := plotter.NewFunction(fn)
	pFn.Color = color
	pFn.Width = vg.Points(1)
	pFn.Samples = 200
	ip.plot.Add(pFn)
}

func (ip *InfoPlotter) AddPoint(x, y float64) {
	point := &plotter.XY{X: x, Y: y}
	if ip.lastPoint != nil {
		l, err := plotter.NewLine(plotter.XYs{*ip.lastPoint, *point})
		if err != nil {
			panic(err)
		}
		if ip.color != nil {
			l.Color = ip.color
		}
		ip.plot.Add(l)
	}
	ip.lastPoint = point
}
