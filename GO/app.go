package main

import (
	"fmt"
	"math"
)

func main() {
	const inflationRate = 7.0
	var investmentAmount float64
	var years float64
	var expectedReturnRate float64
	fmt.Print("Enter the investment amount: ")
	fmt.Scan(&investmentAmount)
	fmt.Print("Enter the years: ")
	fmt.Scan(&years)
	fmt.Print("Enter the expected Return Rate: ")
	fmt.Scan(&expectedReturnRate)

	futureValue := investmentAmount * math.Pow((1+expectedReturnRate/100), years)
	futureRealValue := futureValue / math.Pow((1+inflationRate/100), years)
	fmt.Println(futureValue)
	fmt.Println(futureRealValue)

}
