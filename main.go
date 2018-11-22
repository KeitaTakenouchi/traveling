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
	points []*point
}

func newPath() *path {
	return &path{
		points: make([]*point, 0),
	}
}

func (p *path) addPoint(pt *point) {
	p.points = append(p.points, pt)
}

func (p *path) count() int {
	return len(p.points)
}

func (p *path) distance() float64 {
	sum := 0.0
	totalCount := len(p.points)
	for i := 0; i < totalCount; i++ {
		currPoint := p.points[i]
		nextPoint := p.points[(i+1)%totalCount]

		d := dist(*currPoint, *nextPoint)
		if (i+1)%10 == 0 && !isPrime(currPoint.id) {
			d = d * 1.10
		}
		sum = sum + d
	}
	return sum
}

func isPrime(n int) bool {
	limit := math.Floor(math.Sqrt(float64(n)))
	for i := 2; i < int(limit); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	file, err := os.Open("./cities.csv")
	//file, err := os.Open("./small.csv")
	if err != nil {
		panic("fine not found.")
	}
	reader := csv.NewReader(bufio.NewReader(file))

	path := newPath()
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
		pt := newPoint(i, x, y)
		path.addPoint(pt)
	}

	dist := path.distance()
	fmt.Printf("dist %f", dist)

}
