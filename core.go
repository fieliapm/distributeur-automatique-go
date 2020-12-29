package main

import (
	"fmt"
	"sort"
	"time"
)

type exactPurchaseSolutionKey struct {
	Budget    int
	ItemCount int
}
type exactPurchaseSolutionValue struct {
	SolutionCount int
	Solutions     [][]int
}

type exactPurchaseSolutionMap map[exactPurchaseSolutionKey]exactPurchaseSolutionValue

func (solutionMap exactPurchaseSolutionMap) findExactPurchaseFast(prices []int, budget int, itemCount int) (int, [][]int) {
	exactPurchaseSolutionKey := exactPurchaseSolutionKey{Budget: budget, ItemCount: itemCount}
	exactPurchaseSolutiondValue, ok := solutionMap[exactPurchaseSolutionKey]
	if ok {
		return exactPurchaseSolutiondValue.SolutionCount, exactPurchaseSolutiondValue.Solutions
	} else {
		solutionCount, solutions := solutionMap.findExactPurchaseSlow(prices, budget, itemCount)
		solutionMap[exactPurchaseSolutionKey] = exactPurchaseSolutionValue{SolutionCount: solutionCount, Solutions: solutions}
		return solutionCount, solutions
	}
}

func (solutionMap exactPurchaseSolutionMap) findExactPurchaseSlow(prices []int, budget int, itemCount int) (int, [][]int) {
	if itemCount <= 0 || budget <= 0 {
		if budget == 0 {
			return 1, [][]int{nil}
		} else {
			return 0, nil
		}
	}

	var solutionCount int = 0
	var solutions [][]int

	subSolutionCount, subSolutions := solutionMap.findExactPurchaseFast(prices, budget, itemCount-1)
	solutionCount += subSolutionCount
	solutions = append(solutions, subSolutions...)

	price := prices[itemCount-1]
	subSolutionCount, subSolutions = solutionMap.findExactPurchaseFast(prices, budget-price, itemCount)
	solutionCount += subSolutionCount
	for _, solution := range subSolutions {
		newSolution := make([]int, len(solution)+1)
		copy(newSolution, solution)
		newSolution[len(solution)] = price
		solutions = append(solutions, newSolution)
	}

	return solutionCount, solutions
}

func FindExactPurchaseCache(prices []int, budget int) (int, [][]int) {
	solutionMap := make(exactPurchaseSolutionMap)
	return solutionMap.findExactPurchaseFast(prices, budget, len(prices))
}

func FindExactPurchaseDP(prices []int, budget int) (int, [][]int) {
	dp := make([][]exactPurchaseSolutionValue, len(prices)+1)
	for i, _ := range dp {
		dp[i] = make([]exactPurchaseSolutionValue, budget+1)
	}

	for i, price := range prices {
		dp[i+1][0] = exactPurchaseSolutionValue{SolutionCount: 1, Solutions: [][]int{nil}}
		for j := 1; j <= budget; j++ {
			dp[i+1][j] = dp[i][j]
			if j >= price {
				subSolution := dp[i+1][j-price]
				dp[i+1][j].SolutionCount += subSolution.SolutionCount
				solutions := dp[i+1][j].Solutions
				for _, solution := range subSolution.Solutions {
					newSolution := make([]int, len(solution)+1)
					copy(newSolution, solution)
					newSolution[len(solution)] = price
					solutions = append(solutions, newSolution)
				}
				dp[i+1][j].Solutions = solutions
			}
		}
	}
	result := dp[len(prices)][budget]
	return result.SolutionCount, result.Solutions
}

func FindExactPurchase(prices []int, budget int) int {
	s := time.Now()
	solutionCount, solutions := FindExactPurchaseCache(prices, budget)
	//solutionCount, solutions := FindExactPurchaseDP(prices, budget)
	e := time.Now()
	fmt.Println(float64(e.Sub(s)) / float64(time.Millisecond))

	var c int
	for _, solution := range solutions {
		var v int
		for _, value := range solution {
			v += value
		}
		//fmt.Printf("solution: %v\n", solution)
		if v != budget {
			fmt.Printf("wrong: %v\n", solution)
			fmt.Printf("fq: %d\n", v)
		}

		c++

	}
	if c != solutionCount {
		fmt.Printf("fq2: %d\n", c)
	}

	return solutionCount
}

func main() {
	var prices = []int{24, 29, 62, 37, 33, 22, 109, 38, 32, 75, 57, 30, 132, 19}
	//var prices = []int{2, 3, 5}
	sort.Ints(prices)
	count := FindExactPurchase(prices, 250)
	fmt.Println(count)
}
