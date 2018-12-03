package tsp

import (
	"fmt"
	"math"
	"sort"

	"github.com/fogleman/delaunay"
	"github.com/fogleman/gg"
)

// Algorithm using minimun spannning tree and creatint a tour on it.
func SpannningTreeTourAlgorithm(pool PointPool) *Path {
	edges := spanningTree(pool)

	graph := make(map[*Point][]*Point)
	for _, pt := range pool.Points {
		connected := make([]*Point, 0)
		for _, Edge := range edges {
			if pt == Edge.fst {
				connected = append(connected, Edge.snd)
			} else if pt == Edge.snd {
				connected = append(connected, Edge.fst)
			}
		}
		// sort connected points by their angle.
		sort.Slice(connected, func(i, j int) bool {
			a := connected[i]
			b := connected[j]
			return math.Atan2(a.Y, a.X) < math.Atan2(b.Y, b.X)
		})
		graph[pt] = connected
	}
	tour := make([]*Point, 0)

	startPt := pool.Start
	tour = append(tour, startPt)
	nextOfStartPt := nextPoint(nil, startPt, graph)

	currentPt := nextOfStartPt
	nextPt := nextPoint(startPt, currentPt, graph)
	for !(currentPt == startPt && nextPt == nextOfStartPt) {
		tour = append(tour, currentPt)
		tmp := nextPt
		nextPt = nextPoint(currentPt, nextPt, graph)
		currentPt = tmp
	}

	path := NewPath()
	path.SetStart(pool.Start)
	for _, pt := range tour {
		exist := false
		for _, added := range path.Points {
			if pt == added {
				exist = true
				break
			}
		}
		if !exist {
			path.AddPoint(pt)
		}
	}
	return path
}

func nextPoint(previous, current *Point, graph map[*Point][]*Point) *Point {
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

func nextHalfEdge(e int) int {
	if e%3 == 2 {
		return e - 2
	}
	return e + 1
}

func triangulate(pool PointPool) (*delaunay.Triangulation, error) {
	Points := make([]delaunay.Point, 0)
	for _, pt := range pool.Points {
		p := delaunay.Point{
			X: pt.X,
			Y: pt.Y,
		}
		Points = append(Points, p)
	}
	return delaunay.Triangulate(Points)
}

func spanningTree(pool PointPool) []*Edge {
	fmt.Println("Start triangulation.")
	triangulation, err := triangulate(pool)
	if err != nil {
		panic("triangulation err.")
	}
	//exportTriangulationPNG(triangulation)

	// load edges from the triangulation result.
	edges := make([]*Edge, 0)
	for i, h := range triangulation.Halfedges {
		if i > h {
			p := pool.Points[triangulation.Triangles[i]]
			q := pool.Points[triangulation.Triangles[nextHalfEdge(i)]]
			Edge := NewEdge(p, q)
			edges = append(edges, Edge)
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Distance() < edges[j].Distance()
	})

	fmt.Println("Start spanning.")
	id2gruupId := make([]int, len(pool.Points))

	// init each group ids.
	for i := range id2gruupId {
		id2gruupId[i] = i
	}

	spanningTreeEdges := make([]*Edge, 0)
	for _, Edge := range edges {
		id1, id2 := Edge.fst.ID, Edge.snd.ID

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
		spanningTreeEdges = append(spanningTreeEdges, Edge)
	}

	return spanningTreeEdges
}

func ExportSpanningTreePNG(spanningTreeEdges []*Edge) {
	Points := make([]*Point, 0)
	for _, Edge := range spanningTreeEdges {
		Points = append(Points, Edge.fst)
		Points = append(Points, Edge.snd)
	}

	maxX, maxY := 0.0, 0.0
	for _, pt := range Points {
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

	for _, Edge := range spanningTreeEdges {
		p := Edge.fst
		q := Edge.snd
		ctx.DrawLine(p.X, p.Y, q.X, q.Y)
	}
	ctx.SetRGB(0, 0, 0)
	ctx.Stroke()

	ctx.SavePNG("data/img/spannning.png")
}

func ExportTriangulationPNG(triangulation *delaunay.Triangulation) {
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
