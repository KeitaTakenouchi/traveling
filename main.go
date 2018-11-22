package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

type point struct {
	id int
	x  float64
	y  float64
}

func newPoint(id int, x, y float64) *point {
	return &point{id, x, y}
}

func dist(a, b point) float64 {
	return math.Sqrt((a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y))
}

func main() {
	file, _ := os.Open("./cities.csv")
	reader := csv.NewReader(bufio.NewReader(file))

	points := make([]*point, 0)
	for i := 0; ; i++ {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		x, e1 := strconv.ParseFloat(line[1], 64)
		y, e2 := strconv.ParseFloat(line[2], 64)
		if e1 != nil || e2 != nil {
			continue
		}
		points = append(points, newPoint(i, x, y))
	}

	for _, p := range points {
		fmt.Println(p)
	}

}
