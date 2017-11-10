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

func (p *Polyline) Simplify(tolerance float64) {
	//remove points if they are less than tolerance from both of their neigbhors
	newVertices := make([][]float64, 0)
	newVertices = append(newVertices, p.Vertices[0])
	for j := 1; j < len(p.Vertices)-1; j++ {

		distanceFromNeighborLine := math.Abs((p.Vertices[j+1][1]-newVertices[len(newVertices)-1][1])*p.Vertices[j][0]-(p.Vertices[j+1][0]-newVertices[len(newVertices)-1][0])*p.Vertices[j][1]+p.Vertices[j+1][0]*newVertices[len(newVertices)-1][1]-newVertices[len(newVertices)-1][0]*p.Vertices[j+1][1]) /
			math.Sqrt(math.Pow(p.Vertices[j+1][0]-newVertices[len(newVertices)-1][0], 2)+math.Pow(p.Vertices[j+1][1]-newVertices[len(newVertices)-1][1], 2))

		if distanceFromNeighborLine > tolerance {
			newVertices = append(newVertices, p.Vertices[j])
		}
	}
	fmt.Printf("Simplifying with tolerance %v removed %v of %v points\n", tolerance, len(p.Vertices)-len(newVertices), len(p.Vertices))
	p.Vertices = newVertices
	p.Num = len(newVertices)
}

func (p *Polyline) Center() (float64, float64) {
	xmin, ymin, xmax, ymax := p.BoundingBox()
	return (xmax - xmin) / 2.0, (ymax - ymin) / 2.0
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
	xcenter, ycenter := p.Center()
	return fmt.Sprintf("Object with bounding box of: (%v,%v) to (%v,%v), centerpoint (%v,%v)", xmin, ymin, xmax, ymax, xcenter, ycenter)
}
