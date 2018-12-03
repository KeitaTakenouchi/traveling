package tsp

import (
	"fmt"
	"math"
)

func NearestNextAlgorithm(pool PointPool) *Path {
	totalCount := pool.Size()
	ratio := 0.0

	path := NewPath()
	path.SetStart(pool.Start)

	currentPoint := pool.RemoveAt(0)
	path.AddPoint(currentPoint)
	for !pool.IsEmpty() {
		nextPt := pool.Nearest(currentPoint)
		pool.RemoveById(nextPt.ID)
		path.AddPoint(nextPt)
		currentPoint = nextPt

		// printing info
		r := math.Floor(float64(pool.Size()) / float64(totalCount) * 100)
		if ratio != r {
			ratio = r
			fmt.Printf("*")
		}
	}
	fmt.Println()
	return path
}
