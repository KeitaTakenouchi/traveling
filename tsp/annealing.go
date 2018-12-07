package tsp

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	maxTemperature = 1500.0
	loopLimit      = 1000000
)

func AnnealingAlgorithm(pool PointPool) *Path {
	// set seed of rand.
	rand.Seed(time.Now().UTC().UnixNano())

	// create init state
	initState := NewPath()
	initState.SetStart(pool.Start)
	for _, pt := range pool.Points {
		initState.AddPoint(pt)
	}
	initState.AddPoint(pool.Start)

	adjustmentValue := initState.Distance()

	//start := time.Now().UTC().Unix()

	fmt.Println("start")
	currentState := initState
	ratio := 0
	for i := 1; i <= loopLimit; i++ {
		/*
			if i == 100 {
				now := time.Now().UTC().Unix()
				fmt.Printf("DONE : i=%d,  %d sec\n", i, now-start)
				break
			}
		*/
		r := int(math.Floor(float64(i) / float64(loopLimit) * 1000))
		if ratio != r {
			ratio = r
			if (r % 100) == 0 {
				fmt.Printf("*\n")
				fmt.Printf("DIST : %f\n", currentState.Distance())
			} else {
				adjustmentValue = currentState.Distance()
				fmt.Printf("*")
			}

			//WritePathToFile(currentState, "data/result_annealing.csv")
		}

		temp := maxTemperature * (1.0 - (float64(i) / float64(loopLimit)))

		// create a next candidate.
		nextState := NewPath()
		nextState.SetStart(currentState.Start)
		for _, pt := range currentState.Points {
			nextState.AddPoint(pt)
		}

		i := rand.Intn(len(currentState.Points) - 2)
		k := rand.Intn(len(currentState.Points) - 2 - i)

		a1 := *currentState.Points[i]
		a2 := *currentState.Points[i+1]
		b1 := *currentState.Points[i+1+k]
		b2 := *currentState.Points[i+2+k]

		before := Dist(a1, a2) + Dist(b1, b2) + adjustmentValue
		after := Dist(a1, b1) + Dist(a2, b2) + adjustmentValue
		if after < before {
			// probability = 1
			currentState.Swap(i+1, i+1+k)
		} else if math.Exp((before-after)/temp) > rand.Float64() {
			currentState.Swap(i+1, i+1+k)
		}
	}
	fmt.Println()

	return currentState
}
