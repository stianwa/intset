package intset_test

import (
	"fmt"
	"github.com/stianwa/intset"
)

func ExampleNew() {
	// Create a populated set
	a := intset.New(intset.NegInf(-100), intset.Range(-10,10), intset.PosInf(100))
	
	// Create an empty set
	b := intset.New()

	fmt.Printf("a: %s\nb: %s\n", a, b)
}


func ExampleAll() {
	fmt.Println(intset.New(intset.All()))

	// Same as
	fmt.Println(intset.New().Complement())
}


func ExampleNegInf() {
	// Create an integer set ranging from -∞ to 15
	a := intset.New(intset.NegInf(15))

	// The complement of this set would be a set ranging from 16 to ∞
	fmt.Println(a.Complement())
}


func ExamplePosInf() {
	// Create an integer set ranging from 15 to ∞
	a := intset.New(intset.PosInf(15))

	// The complement of this set would be a set ranging from -∞ to 14
	fmt.Println(a.Complement())
}



func Example_Cardinality() {
	a := intset.New(intset.Range(-100,100), intset.Range(260, 784), intset.Int(900))

	if inf, c := a.Cardinality(); !inf {
		fmt.Printf("cardinality of %s is %d\n", a, c)
	} else {
		fmt.Printf("cardinality of %s is infinite\n", a)
	}
}
