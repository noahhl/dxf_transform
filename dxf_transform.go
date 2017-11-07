package main

import (
	"flag"
	"fmt"
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/entity"
	"math"
)

func main() {

	var inputFile, outputFile string
	var translate, closePoly bool
	var scaleFactor, rotate float64

	flag.StringVar(&inputFile, "in", "", "path to input DXF")
	flag.StringVar(&outputFile, "out", "", "path to write DXF to")
	flag.BoolVar(&translate, "translate", false, "translate DXF to have 0,0 in the lower left corner")
	flag.BoolVar(&closePoly, "close", false, "close polylines")
	flag.Float64Var(&scaleFactor, "scale", 0.0, "scale relative to origin")
	flag.Float64Var(&rotate, "rotate", 0.0, "rotate about origin (prior to translating)")

	flag.Parse()

	drawing, err := dxf.Open(inputFile)
	if err != nil {
		panic(err)
	}

	for i, e := range drawing.Entities() {
		polyline := Polyline{e.(*entity.LwPolyline)}
		if closePoly {
			//if !e.(*entity.LwPolyline).Closed {
			if !polyline.Closed {
				fmt.Printf("Closing entity #%v\n", i)
				polyline.Close()
			}
		}

		if rotate != 0.0 {
			theta := rotate * math.Pi / 180.0
			fmt.Printf("Rotating entity #%v by %v degrees\n", i, rotate)
			vertices := e.(*entity.LwPolyline).Vertices

			for j := range vertices {
				x := vertices[j][0]*math.Cos(theta) + vertices[j][1]*math.Sin(theta)
				y := -1*vertices[j][0]*math.Sin(theta) + vertices[j][1]*math.Cos(theta)
				vertices[j][0] = x
				vertices[j][1] = y
			}

			PrintBoundingBox(vertices)
		}

		if translate {
			//Find global min and max
			var globalMinX = math.MaxFloat64
			var globalMinY = math.MaxFloat64

			for _, e := range drawing.Entities() {
				localMinX, localMinY := FindMins(e.(*entity.LwPolyline).Vertices)
				if localMinX < globalMinX {
					globalMinX = localMinX
				}
				if localMinY < globalMinY {
					globalMinY = localMinY
				}
				PrintBoundingBox(e.(*entity.LwPolyline).Vertices)
			}
			fmt.Printf("Translating entity #%v to have a lower left corner of 0,0\n", i)
			vertices := e.(*entity.LwPolyline).Vertices

			for j := range vertices {
				vertices[j][0] = vertices[j][0] - globalMinX
				vertices[j][1] = vertices[j][1] - globalMinY
			}
			PrintBoundingBox(vertices)
		}
		if scaleFactor > 0.0 {
			fmt.Printf("Scaling entity #%v by %v\n", i, scaleFactor)
			vertices := e.(*entity.LwPolyline).Vertices

			for j := range vertices {
				vertices[j][0] = vertices[j][0] * scaleFactor
				vertices[j][1] = vertices[j][1] * scaleFactor
			}

			PrintBoundingBox(vertices)
		}
	}

	if closePoly || translate || scaleFactor > 0.0 || rotate != 0.0 {
		drawing.SaveAs(outputFile)
	} else {
		fmt.Printf("You didn't ask for any transformations, not saving anything\n")
	}
}

type Polyline struct {
	*entity.LwPolyline
}

func (p *Polyline) BoundingBox() (float64, float64, float64, float64) {
	return 0, 0, 0, 0
}

func PrintBoundingBox(vertices [][]float64) {

	xmin, ymin := FindMins(vertices)
	xmax, ymax := FindMaxs(vertices)
	fmt.Printf("Current bounding box: (%v,%v) to (%v,%v) \n", xmin, ymin, xmax, ymax)
}

func FindMaxs(vertices [][]float64) (float64, float64) {
	var xmax = -math.MaxFloat64
	var ymax = -math.MaxFloat64
	for i := range vertices {
		if vertices[i][0] > xmax {
			xmax = vertices[i][0]
		}
		if vertices[i][1] > ymax {
			ymax = vertices[i][1]
		}
	}

	return xmax, ymax
}

func FindMins(vertices [][]float64) (float64, float64) {
	var xmin = math.MaxFloat64
	var ymin = math.MaxFloat64
	for i := range vertices {
		if vertices[i][0] < xmin {
			xmin = vertices[i][0]
		}
		if vertices[i][1] < ymin {
			ymin = vertices[i][1]
		}
	}

	return xmin, ymin
}
