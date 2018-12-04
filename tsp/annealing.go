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

	currentState := initState
	ratio := 0.0
	for i := 1; i <= loopLimit; i++ {
		r := math.Floor(float64(i) / float64(loopLimit) * 100)
		if ratio != r {
			ratio = r
			fmt.Printf("*")
		}

		temp := temperature(float64(i) / float64(loopLimit))
		nextState := neighbour(currentState)

		if probability(currentState, nextState, temp) > rand.Float64() {
			currentState = nextState
			//fmt.Printf("[%d]\t temp=%f\t, Dist=%f\n", i, temp, currentState.Distance())
		}
	}

	return currentState
}

func temperature(ratio float64) float64 {
	return maxTemperature * (1.0 - ratio)
}

func neighbour(state *Path) *Path {
	next := NewPath()
	next.SetStart(state.Start)
	for _, pt := range state.Points {
		next.AddPoint(pt)
	}

	r1 := rand.Intn(len(state.Points) - 2)
	r2 := rand.Intn(len(state.Points) - 2 - r1)
	next.Swap(r1+1, r1+1+r2) // node 0 is fixed.
	return next
}

func probability(currentState, nextState *Path, temperature float64) float64 {
	if energy(currentState) > energy(nextState) {
		return 1.0
	}
	value := math.Exp((energy(currentState) - energy(nextState)) / temperature)
	return value
}

func energy(state *Path) float64 {
	return state.Distance()
}
