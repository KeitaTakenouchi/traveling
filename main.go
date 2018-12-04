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

func main() {
	//rfile, err := os.Open("data/cities.csv")
	rfile, err := os.Open("data/small.csv")
	defer rfile.Close()
	if err != nil {
		panic("file not found.")
	}
	reader := csv.NewReader(bufio.NewReader(rfile))

	// read points as pool from a file.
	pool := tsp.NewPointPool()
	for i := 0; ; i++ {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		id, e0 := strconv.Atoi(line[0])
		x, e1 := strconv.ParseFloat(line[1], 64)
		y, e2 := strconv.ParseFloat(line[2], 64)
		if e0 != nil || e1 != nil || e2 != nil {
			continue
		}
		pt := tsp.NewPoint(id, x, y)
		pool.AddPoint(pt)
		if pt.ID == 0 {
			pool.SetStart(pt)
		}
	}

	// calculate a path.
	//path := tsp.NearestNextAlgorithm(pool)

	path := tsp.SpannningTreeTourAlgorithm(pool)

	fmt.Printf("dist %f\n", path.Distance())

	tsp.WritePathToFile(path, "data/result.csv")
	tsp.ExportPathPNG(path, "data/img/path.png")
}
