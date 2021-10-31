package intset

import (
	"fmt"
)

type Element struct {
	all    bool
	neginf bool
	posinf bool
	first  int
	last   int
}

// String returns the element set in a human readable form, in
// compliance with the fmt.Stringer interface.
func (e *Element) String() string {
	if e.all {
		return fmt.Sprintf("-%c:%c", 0x221e, 0x221e)
	} else if e.neginf {
		return fmt.Sprintf("-%c:%d", 0x221e, e.first)
	} else if e.posinf {
		return fmt.Sprintf("%d:%c", e.first, 0x221e)
	} else if e.first == e.last {
		return fmt.Sprintf("%d", e.first)
	} else {
		return fmt.Sprintf("%d:%d", e.first, e.last)
	}
}

// Range returns an integer range from a to b.
func Range(a, b int) *Element {
	if b < a {
		a, b = b, a
	}
	return &Element{first: a, last: b}
}

// Int returns a single integer element of the integer n.
func Int(n int) *Element {
	return &Element{first: n, last: n}
}

// All returns a set element spanning from -∞ to ∞ which is the complement of ∅.
func All() *Element {
	return &Element{all: true}
}

// NegInf returns a set element spanning from -∞ to n.
func NegInf(n int) *Element {
	return &Element{first: n, last: n, neginf: true}
}

// PosInf returns a set element spanning from n to ∞.
func PosInf(n int) *Element {
	return &Element{first: n, last: n, posinf: true}
}

// inf Returns true if element set has an infinite flag set
func (e *Element) inf() bool {
	return e.all || e.neginf || e.posinf
}

// isAdjacent Returns true if the two element sets are adjacent.
func (a *Element) isAdjacent(b *Element) bool {
	return  a.all ||
		b.all ||

		(a.neginf && !b.inf() && b.first == a.first + 1) ||
		(b.neginf && !a.inf() && a.first == b.first + 1) ||

		(a.posinf && !b.inf() && a.first - 1 == b.last) ||
		(b.posinf && !a.inf() && b.first - 1 == a.last) ||

		(a.posinf && b.neginf && b.first + 1 == a.first) ||
		(b.posinf && a.neginf && a.first + 1 == b.first) ||

		(!a.inf() && !b.inf() && b.first == a.last + 1) ||
		(!b.inf() && !a.inf() && a.first == b.last + 1)
}

// isOverlapping Returns true if the two element sets are overlapping.
func (a *Element) isOverlapping(b *Element) bool {
	return  a.all  ||
		b.all  ||

		(a.neginf && (b.neginf || b.first <= a.first)) ||
		(b.neginf && (a.neginf || a.first <= b.first)) ||

		(a.posinf && (b.posinf || b.first >= a.first)) ||
		(b.posinf && (a.posinf || a.first >= b.first)) ||

		(b.first <= a.first && b.last >= a.first) ||
		(a.first <= b.first && a.last >= b.first) 
}


// isEqual Returns true if the two element sets are equal.
func (a *Element) isEqual(b *Element) bool {
	return (a.all && b.all) ||
		(a.neginf && b.neginf && a.first == b.first) ||
		(a.posinf && b.posinf && a.first == b.first) ||
		(!a.inf() && !b.inf() && a.first == b.first && a.last == b.last)
}

// isEqual Returns true if the two element sets are equal.
func (a *Element) isSuper(b *Element) bool {
	return  (a.neginf && b.neginf && b.first < a.first) ||
		(a.neginf && !b.inf() && b.last < a.first) ||
		(a.posinf && b.posinf && b.first > a.first) ||
		(a.posinf && !b.inf() && b.first > a.first) ||
		(!a.inf() && !b.inf() && a.first < b.first && a.last > b.last)
}

// isWithin Returns true if element a is within element b.
func (a *Element) isWithin(b *Element) bool {
	return a.first >= b.first && a.last <= b.last
}

// join returns a joined element set when joining two overlapping
// elements a and b. Note! The function does not check for overlap,
// this must be done prior to calling this function.
func (a *Element) join(b *Element) *Element {
	neg, e := minInt(a,b)
	pos, f := maxInt(a,b)

	if neg < 0 && pos < 0 {
		return All()
	} else if neg < 0 {
		return NegInf(f)
	} else if pos < 0 {
		return PosInf(e)
	} else {
		return Range(e,f)
	}
}

// remove returns a list of element sets for removine set b from a.
func (a *Element) remove(b *Element) []*Element {
	var ret []*Element

	if a.isEqual(b) || b.all {
		return ret
	} else if  b.isSuper(a) || ! a.isOverlapping(b) {
		ret = append(ret, a)
		return ret
	}
	
	if a.inf() || b.inf() {
		if a.all && b.neginf {
			ret = append(ret, &Element{first: b.first + 1, posinf: true})
		} else if a.all && b.posinf {
			ret = append(ret, &Element{first: b.first - 1, neginf: true})
		} else if a.posinf && b.posinf {
			ret = append(ret, &Element{first: a.first, last: b.first - 1})
		} else if a.posinf && b.neginf {
			ret = append(ret, &Element{first: b.first + 1, posinf: true})
		} else if a.neginf && b.neginf {
			ret = append(ret, &Element{first: b.first + 1, last: a.first})
		} else if a.neginf && b.posinf {
			ret = append(ret, &Element{first: b.first - 1, neginf: true})
		} else if a.posinf && !b.inf() {
			ret = append(ret, &Element{first: b.last + 1, posinf: true})
		} else if a.neginf && !b.inf() {
			ret = append(ret, &Element{first: b.first - 1, neginf: true})
		} else if !a.inf() && b.posinf {
			ret = append(ret, &Element{first: a.first, last: b.first - 1})
		} else if !a.inf() && b.neginf {
			ret = append(ret, &Element{first: b.first + 1, last: a.last})
		} else if a.all && !b.inf() {
			ret = append(ret, &Element{first: b.first - 1, neginf: true})
			ret = append(ret, &Element{first: b.last  + 1, posinf: true})
		} else {
			fmt.Printf("Here 200 a:%q -b:%q\n", a, b)
		}
	} else if b.isWithin(a) {
		ret = append(ret, &Element{first: a.first, last: b.first - 1})
		ret = append(ret, &Element{first: b.last + 1, last: a.last})
	} else {
		if a.first >= b.first {
			ret = append(ret, &Element{first: b.last + 1, last: a.last})
		} else {
			ret = append(ret, &Element{first: a.first, last: b.first - 1})
		}
	}
	
	return ret
}


// intersect returns a list of element sets from intersecting two
// element sets.
func (a *Element) intersect(b *Element) []*Element {
	var ret []*Element

	if a.isEqual(b) || b.all {
		ret = append(ret, a)
		return ret
	} else if a.isSuper(b) {
		ret = append(ret, b)
		return ret
	} else if b.isSuper(a) {
		ret = append(ret, a)
		return ret
	} else if ! a.isOverlapping(b) {
		return ret
	}
	
	if a.inf() || b.inf() {
		if a.all && b.neginf {
			ret = append(ret, b)
		} else if a.all && b.posinf {
			ret = append(ret, b)
		}  else if a.posinf && b.posinf {
			ret = append(ret, &Element{first: largestOf(a.first,b.first), posinf: true})
		}  else if a.posinf && b.neginf {
			ret = append(ret, &Element{first: a.first, last: b.first})
		}  else if a.neginf && b.neginf {
			ret = append(ret, &Element{first: smallestOf(a.first,b.first), neginf: true})
		}  else if a.neginf && b.posinf {
			ret = append(ret, &Element{first: b.first, last: a.last})
		}  else if a.posinf && !b.inf() {
			ret = append(ret, &Element{first: largestOf(a.first, b.first), last: b.last})
		}  else if a.neginf && !b.inf() {
			ret = append(ret, &Element{first: b.first, last: smallestOf(a.first, b.last)})
		}  else if !a.inf() && b.posinf {
			ret = append(ret, &Element{first: largestOf(a.first, b.first), last: a.last})
		}  else if !a.inf() && b.neginf {
			ret = append(ret, &Element{first: a.first, last: smallestOf(a.last, b.first)})
		}
	} else {
		ret = append(ret, &Element{first: largestOf(a.first, b.first), last: smallestOf(a.last, b.last)})
	}
	
	return ret
}



// Returns the minimum range value of two ranges. Returns two
// integers: 0 or neginf if negative infinitive and the minInt value of a
// and b unless negative infinitive.
func minInt(a, b *Element) (int, int) {
	if a.all || b.all || a.neginf || b.neginf || a.posinf || b.posinf {
		if a.neginf || b.neginf || a.all || b.all {
			return -1, 0
		} else if a.posinf && b.posinf {
			if a.first < b.first {
				return 0, a.first
			} else {
				return 0, b.first
			}
		} else if a.posinf {
			return 0, b.first
		} else {
			return 0, a.first
		}
	}

	if a.first < b.first  {
		return 0, a.first
	}
	return 0, b.first
}

// Returns the upper limit of two ranges. The function returns two
// integers. The first is set to posinf if the maximum point is
// positive infinite. The second is set to the maxumum point unless the
// first is set to posinf.
func maxInt(a, b *Element) (int, int) {
	if a.all || b.all || a.neginf || b.neginf || a.posinf || b.posinf {
		if a.posinf || b.posinf || a.all || b.all {
			return -1, 0
		} else if a.neginf && b.neginf {
			if a.first > b.first {
				return 0, a.first
			} else {
				return 0, b.first
			}
		} else if a.neginf {
			return 0, b.last
		} else {
			return 0, a.last
		}
	}
	
	if a.last > b.last {
		return 0, a.last
	}

	return 0, b.last
}

func smallestOf(a, b int) (int) {
	if a < b {
		return a
	}
	return b
}

func largestOf(a, b int) (int) {
	if a > b {
		return a
	}
	return b
}
