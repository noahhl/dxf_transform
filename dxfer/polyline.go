package dxfer

import (
	"fmt"
	"github.com/yofu/dxf/entity"
	"math"
)

type Polyline struct {
	*entity.LwPolyline
}

func (p *Polyline) Translate(x, y float64) {

	for j := range p.Vertices {
		p.Vertices[j][0] = p.Vertices[j][0] + x
		p.Vertices[j][1] = p.Vertices[j][1] + y
	}
}

func (p *Polyline) Scale(scaleFactor float64) {
	for j := range p.Vertices {
		p.Vertices[j][0] = p.Vertices[j][0] * scaleFactor
		p.Vertices[j][1] = p.Vertices[j][1] * scaleFactor
	}
}

func (p *Polyline) Center() {
	xmin, ymin, xmax, ymax := p.BoundingBox()

	p.Translate((xmin-xmax)/2.0, (ymin-ymax)/2.0)
}

func (p *Polyline) Rotate(theta float64) {

	for j := range p.Vertices {
		x := p.Vertices[j][0]*math.Cos(theta) + p.Vertices[j][1]*math.Sin(theta)
		y := -1*p.Vertices[j][0]*math.Sin(theta) + p.Vertices[j][1]*math.Cos(theta)
		p.Vertices[j][0] = x
		p.Vertices[j][1] = y
	}

}

func (p *Polyline) BoundingBox() (float64, float64, float64, float64) {
	var xmin = math.MaxFloat64
	var ymin = math.MaxFloat64
	var xmax = -math.MaxFloat64
	var ymax = -math.MaxFloat64
	for i := range p.Vertices {
		if p.Vertices[i][0] < xmin {
			xmin = p.Vertices[i][0]
		}
		if p.Vertices[i][1] < ymin {
			ymin = p.Vertices[i][1]
		}
		if p.Vertices[i][0] > xmax {
			xmax = p.Vertices[i][0]
		}
		if p.Vertices[i][1] > ymax {
			ymax = p.Vertices[i][1]
		}
	}

	return xmin, ymin, xmax, ymax
}

func (p *Polyline) Summary() string {
	xmin, ymin, xmax, ymax := p.BoundingBox()
	return fmt.Sprintf("Object with bounding box of: (%v,%v) to (%v,%v)", xmin, ymin, xmax, ymax)
}
