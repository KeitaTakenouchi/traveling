package tsp

import (
	"bufio"
	"fmt"
	"github.com/fogleman/gg"
	"math"
	"os"
)

type Point struct {
	ID int
	X  float64
	Y  float64
}

func NewPoint(ID int, X, y float64) *Point {
	return &Point{ID, X, y}
}

func Dist(a, b Point) float64 {
	return math.Sqrt((a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y))
}

func DistBiased(from, to Point, step int) float64 {
	d := math.Sqrt((from.X-to.X)*(from.X-to.X) + (from.Y-to.Y)*(from.Y-to.Y))
	if (step)%10 == 0 && !IsPrime(from.ID) {
		d = d * 1.1
	}
	return d
}

type Path struct {
	Points []*Point
	Start  *Point
}

func NewPath() *Path {
	return &Path{
		Points: make([]*Point, 0),
	}
}

func (p *Path) SetStart(pt *Point) {
	p.Start = pt
}

func (p *Path) AddPoint(pt *Point) {
	p.Points = append(p.Points, pt)
}

func (p *Path) Count() int {
	return len(p.Points)
}

func (p *Path) Distance() float64 {
	sum := 0.0
	totalCount := len(p.Points)
	for i := 0; i < totalCount; i++ {
		currPoint := p.Points[i]
		nextPoint := p.Points[(i+1)%totalCount]
		sum = sum + DistBiased(*currPoint, *nextPoint, i+1)
	}
	return sum
}

func (p *Path) Swap(i, k int) {
	if !(i <= k && i > 0 && k < p.Count()-1) {
		panic("Invalid index.")
	}

	ps := make([]*Point, 0)
	for j := 0; j < i; j++ {
		ps = append(ps, p.Points[j])
	}
	for j := k; j >= i; j-- {
		ps = append(ps, p.Points[j])
	}
	for j := k + 1; j < p.Count(); j++ {
		ps = append(ps, p.Points[j])
	}
	p.Points = ps
}

type PointPool struct {
	Points []*Point
	Start  *Point
}

func NewPointPool() PointPool {
	return PointPool{
		Points: make([]*Point, 0),
		Start:  nil,
	}
}

func (pp *PointPool) SetStart(pt *Point) {
	pp.Start = pt
}

func (pp *PointPool) AddPoint(pt *Point) {
	pp.Points = append(pp.Points, pt)
}

func (pp *PointPool) RemoveAt(i int) *Point {
	pt := pp.Points[i]
	pp.Points = append(pp.Points[:i], pp.Points[i+1:]...)
	return pt
}

func (pp *PointPool) RemoveById(ID int) *Point {
	targetIndex := -1
	for i, pt := range pp.Points {
		if pt.ID == ID {
			targetIndex = i
			break
		}
	}
	return pp.RemoveAt(targetIndex)
}

func (pp *PointPool) Size() int {
	return len(pp.Points)
}

func (pp *PointPool) IsEmpty() bool {
	return len(pp.Points) == 0
}

func (pp *PointPool) Nearest(target *Point) *Point {
	var nearestPt *Point
	minDist := math.MaxFloat64
	for _, pt := range pp.Points {
		d := Dist(*target, *pt)
		if d < minDist {
			minDist = d
			nearestPt = pt
		}
	}
	return nearestPt
}

type Edge struct {
	fst *Point
	snd *Point
}

func NewEdge(fst, snd *Point) *Edge {
	return &Edge{
		fst: fst,
		snd: snd,
	}
}

func (e *Edge) Distance() float64 {
	return Dist(*e.fst, *e.snd)
}

func IsPrime(n int) bool {
	limit := math.Floor(math.Sqrt(float64(n)))
	for i := 2; i < int(limit); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func ExportPathPNG(path *Path, fileName string) {
	maxX, maxY := 0.0, 0.0
	for _, pt := range path.Points {
		if pt.X > maxX {
			maxX = pt.X
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
	}
	width, height := maxX, maxY
	ctx := gg.NewContext(int(width), int(height))
	ctx.InvertY()
	ctx.DrawRectangle(0, 0, width, height)
	ctx.SetRGB(1, 1, 1)
	ctx.Fill()

	ctx.SetRGB(0.3, 0.3, 0.3)
	for _, pt := range path.Points {
		ctx.DrawPoint(pt.X, pt.Y, 2)
	}
	ctx.Fill()
	ctx.SetRGB(1, 0, 0)
	for _, pt := range path.Points {
		ctx.LineTo(pt.X, pt.Y)
	}
	ctx.LineTo(path.Start.X, path.Start.Y)
	ctx.SetLineWidth(2)
	ctx.Stroke()

	ctx.SavePNG(fileName)
}

func WritePathToFile(path *Path) {
	wfile, err := os.Create("data/result.csv")
	defer wfile.Close()
	if err != nil {
		panic("file creation err.")
	}

	buf := bufio.NewWriter(wfile)
	buf.WriteString("Path\n")
	for _, pt := range path.Points {
		line := fmt.Sprintf("%d\n", pt.ID)
		buf.WriteString(line)
	}
	buf.WriteString("0\n")
	buf.Flush()
}