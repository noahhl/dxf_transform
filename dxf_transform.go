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
	var scaleFactor float64

	flag.StringVar(&inputFile, "in", "", "path to input DXF")
	flag.StringVar(&outputFile, "out", "", "path to write DXF to")
	flag.BoolVar(&translate, "translate", false, "translate DXF to have 0,0 in the lower left corner")
	flag.BoolVar(&closePoly, "close", false, "close polylines")
	flag.Float64Var(&scaleFactor, "scale", 0.0, "scale relative to origin")

	flag.Parse()

	drawing, err := dxf.Open(inputFile)
	if err != nil {
		panic(err)
	}

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
	}
	fmt.Printf("Global xmin %v ymin %v\n", globalMinX, globalMinY)
	if closePoly {
		for i, e := range drawing.Entities() {
			if !e.(*entity.LwPolyline).Closed {
				fmt.Printf("Closing entity #%v\n", i)
				e.(*entity.LwPolyline).Close()
			}
		}
	}

	if translate {
		for i, e := range drawing.Entities() {

			fmt.Printf("Translating entity #%v to have a lower left corner of 0,0\n", i)
			vertices := e.(*entity.LwPolyline).Vertices

			for i := range vertices {
				vertices[i][0] = vertices[i][0] - globalMinX
				vertices[i][1] = vertices[i][1] - globalMinY
			}
			xmin, ymin := FindMins(vertices)
			xmax, ymax := FindMaxs(vertices)
			fmt.Printf("Bounding box: (%v,%v) to (%v,%v) \n", xmin, ymin, xmax, ymax)
		}
	}

	if scaleFactor > 0.0 {
		for i, e := range drawing.Entities() {
			fmt.Printf("Scaling entity #%v by %v\n", i, scaleFactor)
			vertices := e.(*entity.LwPolyline).Vertices

			for i := range vertices {
				vertices[i][0] = vertices[i][0] * scaleFactor
				vertices[i][1] = vertices[i][1] * scaleFactor
			}

			xmin, ymin := FindMins(vertices)
			xmax, ymax := FindMaxs(vertices)
			fmt.Printf("New Bounding box: (%v,%v) to (%v,%v) \n", xmin, ymin, xmax, ymax)
		}
	}

	drawing.SaveAs(outputFile)
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
