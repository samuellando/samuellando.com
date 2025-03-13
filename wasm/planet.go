//go:build js && wasm

package main

import (
	"image/color"
	"math"
	"time"

	"github.com/tdewolff/canvas"
)

type planet struct {
	x      float64
	y      float64
	radPS  float64
	radius float64
	arcs   []*arc
	color  color.Color
	aplha  int

	// Internal kpis
	last time.Time
}

type arc struct {
	angle  float64
	radius float64
}

func ceratePlanet(x, y, radius, radPS float64, nArcs int, color color.Color) *planet {
	nArcsf := float64(nArcs)
	arcs := make([]*arc, nArcs)
	arcGap := math.Pi / nArcsf
	for i := range nArcs {
		angle := float64(i) * arcGap
		arcs[i] = &arc{angle: angle, radius: radius}
	}
	return &planet{
		x:      x,
		y:      y,
		arcs:   arcs,
		radius: radius,
		radPS:  radPS,
		last:   time.Now(),
		color:  color,
		aplha:  13,
	}
}

func (p *planet) draw(ctx *canvas.Context) {
	ctx.Push()
	ctx.SetStrokeColor(p.color)
	r, g, b, _ := p.color.RGBA()
	ctx.SetFill(color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(p.aplha)})
	for _, arc := range p.arcs {
		ctx.DrawPath(p.x, p.y, arc.getPath())
		p.spin()
	}
	ctx.Pop()
}

func (p *planet) spin() {
	elapsed := float64(time.Since(p.last)) / float64(time.Second)
	p.last = time.Now()
	rad := p.radPS * float64(elapsed)
	p.rotate(rad)
}

func (p *planet) rotate(r float64) {
	for _, arc := range p.arcs {
		arc.rotate(r)
	}
}

func (a *arc) rotate(r float64) {
	a.angle += r
	if a.angle >= 2*math.Pi {
		a.angle -= 2 * math.Pi
	}
}

func (a *arc) getPath() *canvas.Path {
	ry := a.radius
	rx := math.Cos(a.angle) * a.radius
	return canvas.Ellipse(rx, ry)
}
