package main

import "fmt"

func main() {
	var revenue float64
	var expenses float64
	var taxRate float64

	fmt.Print("Enter the revenue amount: ")
	fmt.Scan(&revenue)
	fmt.Print("Enter the expenses amount: ")
	fmt.Scan(&expenses)
	fmt.Print("Enter the taxRate: ")
	fmt.Scan(&taxRate)

	ebt := revenue - expenses
	profit := ebt * (1 - taxRate/100)
	ratio := ebt / profit

	fmt.Println(ratio, ebt, profit)
}
