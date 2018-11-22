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

type path struct {
	points []point
}

func newPath() *path {
	return &path{
		points: make([]point, 0),
	}
}

func (p *path) addPoint(pt point) {
	p.points = append(p.points, pt)
}

func (p *path) count() int {
	return len(p.points)
}

func (p *path) distance() float64 {
	sum := 0.0
	for i, currentPoint := range p.points {
		if i == len(p.points)-1 {
			break
		}
		nextPoint := p.points[i+1]
		sum = sum + dist(currentPoint, nextPoint)
	}
	return sum
}

func main() {
	//file, _ := os.Open("./cities.csv")
	file, _ := os.Open("./top100.csv")
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
