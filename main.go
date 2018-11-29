package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/fogleman/delaunay"
	"github.com/fogleman/gg"
	"io"
	"math"
	"os"
	"sort"
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
	start  *point
}

func newPath() *path {
	return &path{
		points: make([]*point, 0),
	}
}

func (p *path) setStart(pt *point) {
	p.start = pt
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

func (p *path) swap(i, k int) {
	if !(i <= k && i > 0 && k < p.count()-1) {
		panic("Invalid index.")
	}

	ps := make([]*point, 0)
	for j := 0; j < i; j++ {
		ps = append(ps, p.points[j])
	}
	for j := k; j >= i; j-- {
		ps = append(ps, p.points[j])
	}
	for j := k + 1; j < p.count(); j++ {
		ps = append(ps, p.points[j])
	}
	p.points = ps
}

type pointPool struct {
	points []*point
	start  *point
}

func newPointPool() pointPool {
	return pointPool{
		points: make([]*point, 0),
		start:  nil,
	}
}

func (pp *pointPool) setStart(pt *point) {
	pp.start = pt
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

type edge struct {
	fst *point
	snd *point
}

func newEdge(fst, snd *point) *edge {
	return &edge{
		fst: fst,
		snd: snd,
	}
}

func (e *edge) distance() float64 {
	return dist(*e.fst, *e.snd)
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

// Algorithm using minimun spannning tree and creatint a tour on it.
func spannningTreeTourAlgorithm(pool pointPool, edges []*edge) *path {
	graph := make(map[*point][]*point)
	for _, pt := range pool.points {
		connected := make([]*point, 0)
		for _, edge := range edges {
			if pt == edge.fst {
				connected = append(connected, edge.snd)
			} else if pt == edge.snd {
				connected = append(connected, edge.fst)
			}
		}
		// sort connected points by their angle.
		sort.Slice(connected, func(i, j int) bool {
			a := connected[i]
			b := connected[j]
			return math.Atan2(a.y, a.x) < math.Atan2(b.y, b.x)
		})
		graph[pt] = connected
		/*
			// dump
			fmt.Print(pt.id, " : ")
			for _, c := range connected {
				fmt.Print(c.id, ", ")
			}
			fmt.Println()
		*/
	}
	tour := make([]*point, 0)

	startPt := pool.start
	tour = append(tour, startPt)
	nextOfStartPt := nextPoint(nil, startPt, graph)

	currentPt := nextOfStartPt
	nextPt := nextPoint(startPt, currentPt, graph)
	for !(currentPt == startPt && nextPt == nextOfStartPt) {
		fmt.Println(currentPt.id, " -> ", nextPt.id)

		tour = append(tour, currentPt)
		tmp := nextPt
		nextPt = nextPoint(currentPt, nextPt, graph)
		currentPt = tmp
	}

	path := newPath()
	path.setStart(pool.start)
	for _, pt := range tour {
		exist := false
		for _, added := range path.points {
			if pt == added {
				exist = true
				break
			}
		}
		if !exist {
			path.addPoint(pt)
		}
	}
	return path
}

func nextPoint(previous, current *point, graph map[*point][]*point) *point {
	connected := graph[current]
	prevIndex := -1
	for i, pt := range connected {
		if pt == previous {
			prevIndex = i
			break
		}
	}
	if previous != nil && prevIndex < 0 {
		panic("index error.")
	}
	return connected[(prevIndex+1)%len(connected)]
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

func nextHalfEdge(e int) int {
	if e%3 == 2 {
		return e - 2
	}
	return e + 1
}

func triangulate(pool pointPool) (*delaunay.Triangulation, error) {
	points := make([]delaunay.Point, 0)
	for _, pt := range pool.points {
		p := delaunay.Point{
			X: pt.x,
			Y: pt.y,
		}
		points = append(points, p)
	}
	return delaunay.Triangulate(points)
}

func exportTriangulationPNG(triangulation *delaunay.Triangulation) {
	maxX, maxY := 0.0, 0.0
	for _, pt := range triangulation.Points {
		if pt.X > maxX {
			maxX = pt.X
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
	}

	ctx := gg.NewContext(int(maxX), int(maxY))
	ctx.InvertY()
	ctx.DrawRectangle(0, 0, maxX, maxY)
	ctx.SetRGB(1, 1, 1)
	ctx.Fill()

	ts := triangulation.Triangles
	hs := triangulation.Halfedges
	for i, h := range hs {
		if i > h {
			p := triangulation.Points[ts[i]]
			q := triangulation.Points[ts[nextHalfEdge(i)]]
			ctx.DrawLine(p.X, p.Y, q.X, q.Y)
		}
	}
	ctx.SetRGB(0, 0, 0)
	ctx.Stroke()

	ctx.SavePNG("data/img/triangle.png")
}

func exportPathPNG(path *path) {
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
	ctx.LineTo(path.start.x, path.start.y)
	ctx.SetLineWidth(2)
	ctx.Stroke()

	ctx.SavePNG("data/img/path.png")
}

func writePathToFile(path *path) {
	wfile, err := os.Create("data/result.csv")
	defer wfile.Close()
	if err != nil {
		panic("file creation err.")
	}

	buf := bufio.NewWriter(wfile)
	buf.WriteString("Path\n")
	for _, pt := range path.points {
		line := fmt.Sprintf("%d\n", pt.id)
		buf.WriteString(line)
	}
	buf.WriteString("0\n")
	buf.Flush()
}

func spanningTree(pool pointPool) []*edge {
	fmt.Println("start triangulation.")
	triangulation, err := triangulate(pool)
	if err != nil {
		panic("triangulation err.")
	}
	//exportTriangulationPNG(triangulation)

	// load edges from the triangulation result.
	edges := make([]*edge, 0)
	for i, h := range triangulation.Halfedges {
		if i > h {
			p := pool.points[triangulation.Triangles[i]]
			q := pool.points[triangulation.Triangles[nextHalfEdge(i)]]
			edge := newEdge(p, q)
			edges = append(edges, edge)
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].distance() < edges[j].distance()
	})

	fmt.Println("start spanning.")
	id2gruupId := make([]int, len(pool.points))

	// init each group ids.
	for i := range id2gruupId {
		id2gruupId[i] = i
	}

	spanningTreeEdges := make([]*edge, 0)
	for _, edge := range edges {
		id1, id2 := edge.fst.id, edge.snd.id

		// skip if both belong to the same group.
		if id2gruupId[id1] == id2gruupId[id2] {
			continue
		}

		minGroupId, maxGroupId := -1, -1
		if id2gruupId[id1] < id2gruupId[id2] {
			minGroupId = id2gruupId[id1]
			maxGroupId = id2gruupId[id2]
		} else {
			minGroupId = id2gruupId[id2]
			maxGroupId = id2gruupId[id1]
		}
		if minGroupId < 0 || maxGroupId < 0 {
			panic("no min/max group ids.")
		}

		for i, groupid := range id2gruupId {
			if groupid == maxGroupId {
				id2gruupId[i] = minGroupId
			}
		}
		spanningTreeEdges = append(spanningTreeEdges, edge)
	}

	return spanningTreeEdges
}

func exportSpanningTreePNG(spanningTreeEdges []*edge) {
	points := make([]*point, 0)
	for _, edge := range spanningTreeEdges {
		points = append(points, edge.fst)
		points = append(points, edge.snd)
	}

	maxX, maxY := 0.0, 0.0
	for _, pt := range points {
		if pt.x > maxX {
			maxX = pt.x
		}
		if pt.y > maxY {
			maxY = pt.y
		}
	}

	ctx := gg.NewContext(int(maxX), int(maxY))
	ctx.InvertY()
	ctx.DrawRectangle(0, 0, maxX, maxY)
	ctx.SetRGB(1, 1, 1)
	ctx.Fill()

	for _, edge := range spanningTreeEdges {
		p := edge.fst
		q := edge.snd
		ctx.DrawLine(p.x, p.y, q.x, q.y)
	}
	ctx.SetRGB(0, 0, 0)
	ctx.Stroke()

	ctx.SavePNG("data/img/spannning.png")
}

func twoOptAlgorithm(path *path) {
	isChanged := true
	for isChanged {
		isChanged = false
		for i := 1; i < path.count()-2; i++ {
			for k := i + 1; k < path.count()-1; k++ {
				a1 := *path.points[i-1]
				a2 := *path.points[i]
				b1 := *path.points[k]
				b2 := *path.points[k+1]

				if dist(a1, a2)+dist(b1, b2) > dist(a1, b1)+dist(a2, b2) {
					path.swap(i, k)
					isChanged = true
				}
			}
		}
	}
}

func main() {
	//rfile, err := os.Open("data/cities.csv")
	rfile, err := os.Open("data/small.csv")
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
		if pt.id == 0 {
			pool.setStart(pt)
		}
	}

	// calculate a path.
	//path := nearestNextAlgorithm(pool)

	edges := spanningTree(pool)
	exportSpanningTreePNG(edges)
	fmt.Println("done spanning.")
	path := spannningTreeTourAlgorithm(pool, edges)
	fmt.Printf("dist %f\n", path.distance())

	twoOptAlgorithm(path)
	fmt.Printf("dist %f\n", path.distance())

	writePathToFile(path)
	exportPathPNG(path)
}
