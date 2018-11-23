package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	//"github.com/fogleman/delaunay"
	"github.com/fogleman/gg"
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

func distBiased(from, to point, step int) float64 {
	d := math.Sqrt((from.x-to.x)*(from.x-to.x) + (from.y-to.y)*(from.y-to.y))
	if (step)%10 == 0 && !isPrime(from.id) {
		d = d * 1.1
	}
	return d
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
		sum = sum + distBiased(*currPoint, *nextPoint, i+1)
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

func nearestNextAlgorithm(pool pointPool) *path {
	totalCount := pool.size()
	ratio := 0.0

	path := newPath()

	currentPoint := pool.removeAt(0)
	path.addPoint(currentPoint)
	for !pool.isEmpty() {
		nextPt := pool.nearest(currentPoint)
		pool.removeById(nextPt.id)
		path.addPoint(nextPt)
		currentPoint = nextPt

		// printing info
		r := math.Floor(float64(pool.size()) / float64(totalCount) * 100)
		if ratio != r {
			ratio = r
			fmt.Printf("*")
		}
	}
	fmt.Println()
	return path
}

func main() {
	rfile, err := os.Open("data/cities.csv")
	//rfile, err := os.Open("data/small.csv")
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
	fmt.Printf("dist %f\n", dist)

	// write the result as csv.
	wfile, err := os.Create("data/result.csv")
	defer wfile.Close()

	buf := bufio.NewWriter(wfile)
	buf.WriteString("Path\n")
	for _, pt := range path.points {
		line := fmt.Sprintf("%d\n", pt.id)
		buf.WriteString(line)
	}
	buf.WriteString("0\n")
	buf.Flush()

	// output png file.
	maxX, maxY := 0.0, 0.0
	for _, pt := range path.points {
		if pt.x > maxX {
			maxX = pt.x
		}
		if pt.y > maxY {
			maxY = pt.y
		}
	}
	width, height := maxX, maxY
	ctx := gg.NewContext(int(width), int(height))
	ctx.InvertY()
	ctx.DrawRectangle(0, 0, width, height)
	ctx.SetRGB(1, 1, 1)
	ctx.Fill()

	ctx.SetRGB(0.3, 0.3, 0.3)
	for _, pt := range path.points {
		ctx.DrawPoint(pt.x, pt.y, 2)
	}
	ctx.Fill()
	ctx.SetRGB(1, 0, 0)
	for _, pt := range path.points {
		ctx.LineTo(pt.x, pt.y)
	}
	ctx.SetLineWidth(2)
	ctx.Stroke()

	ctx.SavePNG("data/img/out.png")

}
