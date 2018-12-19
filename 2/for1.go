package main

import (
	"fmt"
	"strconv"
)

func main() {
	//for numer i
	i := 0
	for {
		if i > 20 {
			break
		}
		i++
		fmt.Println(i)
	}
	fmt.Println("Wyjście z pętli i: " + strconv.Itoa(i))

	//for numer 2
	var i1 int = 0
	for i1 = 0; i1 < 20; i1++ {
		fmt.Println(i1)
	}

	fmt.Println("Wyjście z pętli i1: " + strconv.Itoa(i1))

	wynik3 := 0
	for i2 := 0; i2 < 30; i2++ {
		wynik3 = i2
	}

	fmt.Printf("%d\n", wynik3)
	fmt.Printf("Dupa Dupa")
}
