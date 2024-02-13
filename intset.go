// Package intset implements set theory methods for sets in the ℤ
// domain. The operations can handle sets in the range -∞:∞, but the
// minimum and maximum values of the integers are limited by the
// underlying 32-bit or a 64-bit machine platform.
//
//	package main
//
//	import (
//	        "github.com/stianwa/intset"
//	        "fmt"
//	)
//
//	func main() {
//	    a := intset.New(intset.Range(-300, -30), intset.NegInf(-500), intset.PosInf(500))
//
//	    b := intset.New(intset.Range(-47, 23))
//
//	    fmt.Printf("%s union %s = %s\n", a, b, a.Union(b))
//	    fmt.Printf("%s intersect %s = %s\n", a, b, a.Intersect(b))
//	    fmt.Printf("complement of %s = %s\n", a, a.Complement())
//
//	    if c, inf := b.Cardinality(); !inf {
//	        fmt.Printf("cardinality of %s is: %d\n", b, c)
//	    } else {
//	        fmt.Println("cardinality of %s is infinite")
//	    }
//	}
package intset

import (
	"fmt"
	"strings"
)

// IntSet holds a slice of element which makes a set.
type IntSet struct {
	elements []*Element
}

// New returns a new set. Any Range sets passed to New, will be added
// to the set.
func New(elements ...*Element) *IntSet {
	n := &IntSet{}
	n.AddElements(elements...)

	return n
}

// AddInts adds integers to a set.
func (a *IntSet) AddInts(numbers ...int) {
	for _, n := range numbers {
		a.insertElement(Int(n))
	}
}

// AddPosInf adds a range from n to ∞ to the set.
func (a *IntSet) AddPosInf(n int) {
	a.insertElement(PosInf(n))

}

// AddNegInf adds a range from -∞ to n to the set.
func (a *IntSet) AddNegInf(n int) {
	a.insertElement(NegInf(n))
}

// RemoveInts removes integers from a set.
func (a *IntSet) RemoveInts(numbers ...int) {
	for _, n := range numbers {
		a.removeElement(Int(n))
	}
}

// AddElements adds element types of all kinds to a set.
func (a *IntSet) AddElements(elements ...*Element) {
	for _, r := range elements {
		a.insertElement(r)
	}
}

// RemoveElements removes elements from a set.
func (a *IntSet) RemoveElements(elements ...*Element) {
	for _, r := range elements {
		a.removeElement(r)
	}
}

// insertRange inserts a single Range to a set.
func (a *IntSet) insertElement(r *Element) {
	if len(a.elements) == 0 {
		a.elements = append(a.elements, r)
		return
	} else if len(a.elements) == 1 && a.elements[0].all { // special case for 'all'
		return
	}

	var newList []*Element
	var prev *Element

	inserted := false
	for i, e := range a.elements {
		if inserted {
			// overlap check mode
			if prev.isOverlapping(e) {
				a.elements[i-1] = prev.join(e)
				continue
			}
		} else {
			// insert mode
			if e.isOverlapping(r) || e.isAdjacent(r) {
				a.elements[i] = e.join(r)

				inserted = true
			} else if r.first < e.first {
				newList = append(newList, r)
				inserted = true
			}
		}
		prev = a.elements[i]
		newList = append(newList, a.elements[i])
	}
	if !inserted {
		newList = append(newList, r)
	}

	a.elements = newList
}

// optimize range sets
func (a *IntSet) optimize() {
	// keep joining ranges until no more ranges can be joined
	for {
		n2 := &IntSet{}
		for _, r := range a.elements {
			n2.insertElement(r)
		}

		if n2.Equal(a) {
			break
		}
		a.elements = n2.elements
	}
}

// removeElement removes a single element from a set.
func (a *IntSet) removeElement(r *Element) {
	var newList []*Element
	for _, e := range a.elements {
		for _, n := range e.remove(r) {
			newList = append(newList, n)
		}
	}
	a.elements = newList
}

// String returns the set in a human readable form, in compliance with
// the fmt.Stringer interface.
func (a *IntSet) String() string {
	if len(a.elements) == 0 {
		return fmt.Sprintf("{%c}", 0x2205)
	}

	var ents []string
	for _, r := range a.elements {
		ents = append(ents, fmt.Sprintf("%s", r))
	}

	return fmt.Sprintf("{%s}", strings.Join(ents, ", "))
}

// HasInt returns true if the integer is part of the set.
func (a *IntSet) HasInt(m int) bool {
	if len(a.elements) == 0 {
		return false
	}

	for _, r := range a.elements {
		if r.all ||
			(r.neginf && m <= r.first) ||
			(r.posinf && m >= r.first) ||
			(m >= r.first && m <= r.last) {
			return true
		}
	}

	return false
}

// Cardinality returns an unsigned integer holding the cardinality and
// an infinite boolean. If the infinite boolean is true, the
// cardinality of the set can not be held by the unsigned integer, and
// the value must be discarded.
func (a *IntSet) Cardinality() (uint, bool) {
	var cardinality uint

	for _, r := range a.elements {
		if r.inf() {
			return 0, true
		}
		if r.last > 0 && r.first < 0 {
			if uintAddOverflow(&cardinality, uint(r.last)) ||
				uintAddOverflow(&cardinality, uint(r.first*-1)) ||
				uintAddOverflow(&cardinality, 1) {
				return 0, true
			}
		} else {
			if uintAddOverflow(&cardinality, uint(r.last-r.first+1)) {
				return 0, true
			}
		}
	}

	return cardinality, false
}

// add uint, returns true if overflow
func uintAddOverflow(p *uint, n uint) bool {
	save := *p
	*p = *p + n
	if *p < save {
		return true
	}

	return false
}

// Complement returns a∁.
func (a *IntSet) Complement() *IntSet {
	n := &IntSet{}

	if len(a.elements) == 0 {
		n.elements = append(n.elements, &Element{all: true})
		return n
	} else if len(a.elements) == 1 {
		l := a.elements[0]
		if l.all {
			a.elements = nil
		} else if l.neginf {
			n.elements = append(n.elements, &Element{first: l.first + 1, posinf: true})
		} else if l.posinf {
			n.elements = append(n.elements, &Element{first: l.first - 1, neginf: true})
		} else {
			n.elements = append(n.elements, &Element{first: l.first - 1, neginf: true})
			n.elements = append(n.elements, &Element{first: l.last + 1, posinf: true})
		}
		return n
	}

	last := len(a.elements) - 2
	prev := a.elements[0]
	for i, e := range a.elements[1:] {
		if i == 0 {
			if !prev.inf() {
				n.elements = append(n.elements, &Element{first: prev.first - 1, neginf: true})
				n.elements = append(n.elements, &Element{first: prev.last + 1, last: e.first - 1})
			} else {
				n.elements = append(n.elements, &Element{first: prev.first + 1, last: e.first - 1})
			}
		} else if i == last && !e.inf() {
			n.elements = append(n.elements, &Element{first: prev.last + 1, last: e.last - 1})
			n.elements = append(n.elements, &Element{first: e.last + 1, posinf: true})
		} else {
			n.elements = append(n.elements, &Element{first: prev.last + 1, last: e.first - 1})
		}
		prev = e
	}

	return n
}

// Union returns a ∪ b.
func (a *IntSet) Union(b *IntSet) *IntSet {
	n := &IntSet{}

	for _, r := range a.elements {
		n.insertElement(r)
	}
	for _, r := range b.elements {
		n.insertElement(r)
	}

	n.optimize()

	return n
}

// Intersect returns a ∩ b.
func (a *IntSet) Intersect(b *IntSet) *IntSet {
	n := &IntSet{}
	for _, ar := range a.elements {
		for _, br := range b.elements {
			for _, e := range ar.intersect(br) {
				n.AddElements(e)
			}
		}
	}
	n.optimize()
	return n
}

// Difference returns a - b.
func (a *IntSet) Difference(b *IntSet) *IntSet {
	n := a.Copy()
	n.RemoveElements(b.elements...)
	n.optimize()

	return n
}

// Xor returns a ⊻ b.
func (a *IntSet) Xor(b *IntSet) *IntSet {
	return a.Union(b).Difference(a.Intersect(b))
}

// Copy returns a copy hf the set.
func (a *IntSet) Copy() *IntSet {
	n := &IntSet{}
	n.AddElements(a.elements...)

	return n
}

// Equal returns true if the two sets are equal.
func (a *IntSet) Equal(b *IntSet) bool {
	if len(a.elements) != len(b.elements) {
		return false
	}

	for i, r := range a.elements {
		if !r.isEqual(b.elements[i]) {
			return false
		}
	}

	return true
}

// IsSubsetOf returns true if a ⊆ b.
func (a *IntSet) IsSubsetOf(b *IntSet) bool {
	return a.Union(b).Equal(b)
}

// IsProperSubsetOf returns true if a ⊊ b.
func (a *IntSet) IsProperSubsetOf(b *IntSet) bool {
	return a.IsSubsetOf(b) && !a.Equal(b)
}
