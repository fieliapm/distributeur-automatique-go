package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.rayark.com/fieliapm/distributeur_automatique/pkg/api"
	"gitlab.rayark.com/fieliapm/distributeur_automatique/pkg/core"
)

var (
	// revision info
	version    string = ""
	commitHash string = ""
	buildDate  string = ""
	// platform info
	goOS   string = ""
	goArch string = ""
)

func showInfo() {
	fmt.Fprintln(os.Stderr, "distributeur automatique")
	fmt.Fprintln(os.Stderr, "written by fieliapm\n")
	fmt.Fprintf(os.Stderr, "Version: %s\nCommitHash: %s\nBuildDate: %s\n\n", version, commitHash, buildDate)
	fmt.Fprintf(os.Stderr, "OS: %s\tArchitecture: %s\n\n", goOS, goArch)
}

func findExactPurchase(prices []int, budget int) {
	startTime := time.Now()

	solutionCount, solutions := core.FindExactPurchaseCache(prices, budget)
	//solutionCount, solutions := core.FindExactPurchaseDP(prices, budget)

	endTime := time.Now()
	fmt.Printf("execution time: %f ms\n", float64(endTime.Sub(startTime))/float64(time.Millisecond))

	fmt.Printf("solution count: %d\n", solutionCount)

	var count int
	for _, solution := range solutions {
		var v int
		for _, value := range solution {
			v += value
		}
		fmt.Printf("solution: %v\n", solution)

		if v != budget {
			fmt.Printf("wrong solution: %v sum:\n", solution, v)
		}
		count++
	}
	if count != solutionCount {
		fmt.Printf("wrong solution count: %d\n", count)
	}
}

func input() (prices []int, budget int, err error) {
	fmt.Fprintln(os.Stderr, "avaliable item count:")
	var priceCount int
	_, err = fmt.Scanf("%d", &priceCount)
	if err != nil {
		return
	}

	fmt.Fprintln(os.Stderr, "avaliable item prices:")
	for i := 0; i < priceCount; i++ {
		var price int
		_, err = fmt.Scanf("%d", &price)
		if err != nil {
			return
		}
		prices = append(prices, price)
	}

	fmt.Fprintln(os.Stderr, "your budget:")
	_, err = fmt.Scanf("%d", &budget)
	if err != nil {
		return
	}

	return
}

func cmd() {
	prices, budget, err := input()
	if err != nil {
		panic(err)
	}

	findExactPurchase(prices, budget)
}

func main() {
	showInfo()

	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "daemon" {
			api.RunServer(8000)
		}
	}
	cmd()
}
