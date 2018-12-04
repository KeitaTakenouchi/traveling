package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/KeitaTakenouchi/traveling/tsp"
)

func twoOptAlgorithm(path *tsp.Path) {
	loop := 0
	isChanged := true
	for isChanged {
		isChanged = false
		for i := 1; i < path.Count()-2; i++ {
			for k := i + 1; k < path.Count()-1; k++ {
				a1 := *path.Points[i-1]
				a2 := *path.Points[i]
				b1 := *path.Points[k]
				b2 := *path.Points[k+1]

				if tsp.Dist(a1, a2)+tsp.Dist(b1, b2) > tsp.Dist(a1, b1)+tsp.Dist(a2, b2) {
					path.Swap(i, k)
					isChanged = true
				}
			}
		}
		loop++
		fmt.Printf("Dist %f\n", path.Distance())
		tsp.WritePathToFile(path, "data/result_2opt.csv")
	}
	fmt.Println()
	fmt.Println("2 opt done. loop = ", loop)
}

func main() {
	//rfile, err := os.Open("data/cities.csv")
	rfile, err := os.Open("data/result.csv")
	defer rfile.Close()
	if err != nil {
		panic("file not found.")
	}
	reader := csv.NewReader(bufio.NewReader(rfile))

	// read Path from a file.
	path := tsp.NewPath()
	for i := 0; ; i++ {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if len(line) != 3 {
			continue
		}
		id, e0 := strconv.Atoi(line[0])
		x, e1 := strconv.ParseFloat(line[1], 64)
		y, e2 := strconv.ParseFloat(line[2], 64)
		if e0 != nil || e1 != nil || e2 != nil {
			continue
		}
		pt := tsp.NewPoint(id, x, y)
		path.AddPoint(pt)
		if pt.ID == 0 {
			path.SetStart(pt)
		}
	}

	twoOptAlgorithm(path)
	tsp.WritePathToFile(path, "data/result_2opt.csv")
	tsp.ExportPathPNG(path, "data/img/path_2opt.png")
}
