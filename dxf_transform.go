package main

import (
	"flag"
	"fmt"
	"github.com/noahhl/dxf_transform/dxfer"
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/entity"
	"math"
)

func main() {

	var inputFile, outputFile string
	var translate, closePoly bool
	var scaleFactor, rotate, tolerance float64

	flag.StringVar(&inputFile, "in", "", "path to input DXF")
	flag.StringVar(&outputFile, "out", "", "path to write DXF to")
	flag.BoolVar(&translate, "translate", false, "translate DXF to have 0,0 in the lower left corner")
	flag.BoolVar(&closePoly, "close", false, "close polylines")
	flag.Float64Var(&scaleFactor, "scale", 0.0, "scale relative to origin")
	flag.Float64Var(&rotate, "rotate", 0.0, "rotate about origin (prior to translating)")
	flag.Float64Var(&tolerance, "tolerance", 0.0, "simplification tolerance (in output units/scale)")

	flag.Parse()

	drawing, err := dxf.Open(inputFile)
	if err != nil {
		panic(err)
	}

	for i, e := range drawing.Entities() {
		polyline := dxfer.Polyline{e.(*entity.LwPolyline)}
		fmt.Printf("\n%v\n", polyline.Summary())

		if closePoly {
			if !polyline.Closed {
				fmt.Printf("Closing entity #%v\n", i)
				polyline.Close()
			}
		}

		if rotate != 0.0 {
			fmt.Printf("Rotating entity #%v by %v degrees\n", i, rotate)
			polyline.Rotate(rotate * math.Pi / 180.0)
			fmt.Printf("%v\n", polyline.Summary())
		}

		if translate {
			//Find global min and max
			var globalMinX = math.MaxFloat64
			var globalMinY = math.MaxFloat64

			for _, e := range drawing.Entities() {
				p := dxfer.Polyline{e.(*entity.LwPolyline)}
				localMinX, localMinY, _, _ := p.BoundingBox()
				if localMinX < globalMinX {
					globalMinX = localMinX
				}
				if localMinY < globalMinY {
					globalMinY = localMinY
				}
			}
			fmt.Printf("Translating entity #%v to have a lower left corner of 0,0\n", i)
			polyline.Translate(-1*globalMinX, -1*globalMinY)
			fmt.Printf("%v\n", polyline.Summary())
		}
		if scaleFactor > 0.0 {
			fmt.Printf("Scaling entity #%v by %v\n", i, scaleFactor)
			polyline.Scale(scaleFactor)
			fmt.Printf("%v\n", polyline.Summary())
		}

		if tolerance > 0.0 {
			polyline.Simplify(tolerance)
		}
	}

	if closePoly || translate || scaleFactor > 0.0 || rotate != 0.0 {
		drawing.SaveAs(outputFile)
		fmt.Printf("\nSaved as %v\n", outputFile)
	} else {
		fmt.Printf("You didn't ask for any transformations, not saving anything\n")
	}
}
