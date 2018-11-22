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

type pointPool struct {
	points []*point
}

func (pp *pointPool) addPoint(pt *point) {
	pp.points = append(pp.points, pt)
}

func (pp *pointPool) removeAt(i int) *point {
	pt := pp.points[i]
	pp.points = append(pp.points[:i], pp.points[i+1:]...)
	return pt
}

func (pp *pointPool) removeById(id int) *point {
	targetIndex := -1
	for i, pt := range pp.points {
		if pt.id == id {
			targetIndex = i
			break
		}
	}
	return pp.removeAt(targetIndex)
}

func (pp *pointPool) size() int {
	return len(pp.points)
}

func (pp *pointPool) isEmpty() bool {
	return len(pp.points) == 0
}

func (pp *pointPool) nearest(target *point) *point {
	var nearestPt *point
	minDist := math.MaxFloat64
	for _, pt := range pp.points {
		d := dist(*target, *pt)
		if d < minDist {
			minDist = d
			nearestPt = pt
		}
	}
	return nearestPt
}

func newPointPool() pointPool {
	return pointPool{
		points: make([]*point, 0),
	}
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
	rfile, err := os.Open("./cities.csv")
	//rfile, err := os.Open("./small.csv")
	defer rfile.Close()
	if err != nil {
		panic("fine not found.")
	}
	reader := csv.NewReader(bufio.NewReader(rfile))

	// read points as pool from a file.
	pool := newPointPool()
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
		pt := newPoint(id, x, y)
		pool.addPoint(pt)
	}

	// calculate a path.
	path := nearestNextAlgorithm(pool)
	dist := path.distance()
	fmt.Printf("dist %f", dist)

	// write the result as csv.
	wfile, err := os.Create("./result.csv")
	defer wfile.Close()

	buf := bufio.NewWriter(wfile)
	buf.WriteString("Path, X, Y\n")
	for _, pt := range path.points {
		line := fmt.Sprintf("%d, %f, %f\n", pt.id, pt.x, pt.y)
		buf.WriteString(line)
	}

	start := path.points[0]
	line := fmt.Sprintf("%d, %f, %f\n", start.id, start.x, start.y)
	buf.WriteString(line)

	buf.Flush()
}

func nearestNextAlgorithm(pool pointPool) *path {
	path := newPath()

	currentPoint := pool.removeAt(0)
	path.addPoint(currentPoint)
	for !pool.isEmpty() {
		nextPt := pool.nearest(currentPoint)
		pool.removeById(nextPt.id)
		path.addPoint(nextPt)
		currentPoint = nextPt
	}
	return path
}
