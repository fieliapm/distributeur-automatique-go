package main

import (
	"fmt"
	"sort"
)

var count int

func FindExactPurchase(prices []int, budget int, history []int) {
	if budget <= 0 {
		if budget == 0 {
			fmt.Println(history)
			count++
		}
		return
	}
	for i, price := range prices {
		remainBudget := budget - price
		FindExactPurchase(prices[i:], remainBudget, append(history, price))
	}
}

func main() {
	var prices = []int{24, 29, 62, 37, 33, 22, 109, 38, 32, 75, 57, 30, 132, 19}
	sort.Ints(prices)
	count = 0
	FindExactPurchase(prices, 250, nil)
	fmt.Println(count)
}
