package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/noahhl/dxf_transform/dxfer"
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/entity"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type TransformationSpec struct {
	Plate string
	File  string
	X     float64
	Y     float64
	Rot   float64
}

var polylines []dxfer.Polyline

func main() {

	var inputFile string

	flag.StringVar(&inputFile, "in", "", "configuration file of things to be nested")

	flag.Parse()

	transformsFile, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	r := csv.NewReader(strings.NewReader(string(transformsFile)))

	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	var plates = make(map[string][]TransformationSpec, 0)
	for _, r := range records {
		if r[0] == "plate" {
			continue
		}

		plate := r[0]
		file := strings.Replace(r[1], ".stl", ".dxf", -1)
		x, _ := strconv.ParseFloat(r[2], 64)
		y, _ := strconv.ParseFloat(r[3], 64)
		rot, _ := strconv.ParseFloat(r[4], 64)
		t := TransformationSpec{plate, file, x, y, rot}
		if _, ok := plates[plate]; !ok {
			plates[plate] = make([]TransformationSpec, 0)
		}
		plates[plate] = append(plates[plate], t)
	}

	for plate, transforms := range plates {

		combinedDraw := dxf.NewDrawing()
		for _, f := range transforms {
			drawing, err := dxf.Open(f.File)
			if err != nil {
				panic(err)
			}

			var drawingXMin = math.MaxFloat64
			var drawingYMin = math.MaxFloat64
			var drawingXMax = -math.MaxFloat64
			var drawingYMax = -math.MaxFloat64
			for _, e := range drawing.Entities() {
				p := dxfer.Polyline{e.(*entity.LwPolyline)}
				lineMinX, lineMinY, lineMaxX, lineMaxY := p.BoundingBox()
				if lineMinX < drawingXMin {
					drawingXMin = lineMinX
				}
				if lineMinY < drawingYMin {
					drawingYMin = lineMinY
				}
				if lineMaxX > drawingXMax {
					drawingXMax = lineMaxX
				}
				if lineMaxY > drawingYMax {
					drawingYMax = lineMaxY
				}
			}

			for _, e := range drawing.Entities() {
				p := dxfer.Polyline{e.(*entity.LwPolyline)}
				p.Translate((drawingXMin-drawingXMax)/2.0, (drawingYMin-drawingYMax)/2.0)
				p.Rotate(-1.0 * f.Rot * math.Pi / 180.0)
				p.Translate(f.X, f.Y)
				combinedDraw.AddEntity(p)
				polylines = append(polylines, p)
			}
		}

		var globalMinX = math.MaxFloat64
		var globalMinY = math.MaxFloat64

		for _, p := range polylines {
			localMinX, localMinY, _, _ := p.BoundingBox()
			if localMinX < globalMinX {
				globalMinX = localMinX
			}
			if localMinY < globalMinY {
				globalMinY = localMinY
			}
		}
		for _, p := range polylines {
			p.Translate(-1*globalMinX, -1*globalMinY)
		}

		combinedDraw.SaveAs(fmt.Sprintf("%v.dxf", plate))
	}

}
