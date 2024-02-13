package intset

import (
	"fmt"
)

// Element stores integers or ranges used in IntSet.
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
	}
	return fmt.Sprintf("%d:%d", e.first, e.last)
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
func (e *Element) isAdjacent(o *Element) bool {
	return e.all ||
		o.all ||

		(e.neginf && !o.inf() && o.first == e.first+1) ||
		(o.neginf && !e.inf() && e.first == o.first+1) ||

		(e.posinf && !o.inf() && e.first-1 == o.last) ||
		(o.posinf && !e.inf() && o.first-1 == e.last) ||

		(e.posinf && o.neginf && o.first+1 == e.first) ||
		(o.posinf && e.neginf && e.first+1 == o.first) ||

		(!e.inf() && !o.inf() && o.first == e.last+1) ||
		(!o.inf() && !e.inf() && e.first == o.last+1)
}

// isOverlapping Returns true if the two element sets are overlapping.
func (e *Element) isOverlapping(o *Element) bool {
	return e.all ||
		o.all ||

		(e.neginf && (o.neginf || o.first <= e.first)) ||
		(o.neginf && (e.neginf || e.first <= o.first)) ||

		(e.posinf && (o.posinf || o.first >= e.first)) ||
		(o.posinf && (e.posinf || e.first >= o.first)) ||

		(o.first <= e.first && o.last >= e.first) ||
		(e.first <= o.first && e.last >= o.first)
}

// isEqual Returns true if the two element sets are equal.
func (e *Element) isEqual(o *Element) bool {
	return (e.all && o.all) ||
		(e.neginf && o.neginf && e.first == o.first) ||
		(e.posinf && o.posinf && e.first == o.first) ||
		(!e.inf() && !o.inf() && e.first == o.first && e.last == o.last)
}

// isEqual Returns true if the two element sets are equal.
func (e *Element) isSuper(o *Element) bool {
	return (e.neginf && o.neginf && o.first < e.first) ||
		(e.neginf && !o.inf() && o.last < e.first) ||
		(e.posinf && o.posinf && o.first > e.first) ||
		(e.posinf && !o.inf() && o.first > e.first) ||
		(!e.inf() && !o.inf() && e.first < o.first && e.last > o.last)
}

// isWithin Returns true if element a is within element o.
func (e *Element) isWithin(o *Element) bool {
	return e.first >= o.first && e.last <= o.last
}

// join returns a joined element set when joining two overlapping
// elements e and o. Note! The function does not check for overlap,
// this must be done prior to calling this function.
func (e *Element) join(o *Element) *Element {
	neg, min := minInt(e, o)
	pos, max := maxInt(e, o)

	if neg < 0 && pos < 0 {
		return All()
	} else if neg < 0 {
		return NegInf(max)
	} else if pos < 0 {
		return PosInf(min)
	}
	return Range(min, max)
}

// remove returns a list of element sets for removine set b from e.
func (e *Element) remove(o *Element) []*Element {
	var ret []*Element

	if e.isEqual(o) || o.all {
		return ret
	} else if o.isSuper(e) || !e.isOverlapping(o) {
		ret = append(ret, e)
		return ret
	}

	if e.inf() || o.inf() {
		if e.all && o.neginf {
			ret = append(ret, &Element{first: o.first + 1, posinf: true})
		} else if e.all && o.posinf {
			ret = append(ret, &Element{first: o.first - 1, neginf: true})
		} else if e.posinf && o.posinf {
			ret = append(ret, &Element{first: e.first, last: o.first - 1})
		} else if e.posinf && o.neginf {
			ret = append(ret, &Element{first: o.first + 1, posinf: true})
		} else if e.neginf && o.neginf {
			ret = append(ret, &Element{first: o.first + 1, last: e.first})
		} else if e.neginf && o.posinf {
			ret = append(ret, &Element{first: o.first - 1, neginf: true})
		} else if e.posinf && !o.inf() {
			ret = append(ret, &Element{first: o.last + 1, posinf: true})
		} else if e.neginf && !o.inf() {
			ret = append(ret, &Element{first: o.first - 1, neginf: true})
		} else if !e.inf() && o.posinf {
			ret = append(ret, &Element{first: e.first, last: o.first - 1})
		} else if !e.inf() && o.neginf {
			ret = append(ret, &Element{first: o.first + 1, last: e.last})
		} else if e.all && !o.inf() {
			ret = append(ret, &Element{first: o.first - 1, neginf: true})
			ret = append(ret, &Element{first: o.last + 1, posinf: true})
		} else {
			fmt.Printf("Here 200 a:%q -b:%q\n", e, o)
		}
	} else if o.isWithin(e) {
		ret = append(ret, &Element{first: e.first, last: o.first - 1})
		ret = append(ret, &Element{first: o.last + 1, last: e.last})
	} else {
		if e.first >= o.first {
			ret = append(ret, &Element{first: o.last + 1, last: e.last})
		} else {
			ret = append(ret, &Element{first: e.first, last: o.first - 1})
		}
	}

	return ret
}

// intersect returns a list of element sets from intersecting two
// element sets.
func (e *Element) intersect(o *Element) []*Element {
	var ret []*Element

	if e.isEqual(o) || o.all {
		ret = append(ret, e)
		return ret
	} else if e.isSuper(o) {
		ret = append(ret, o)
		return ret
	} else if o.isSuper(e) {
		ret = append(ret, e)
		return ret
	} else if !e.isOverlapping(o) {
		return ret
	}

	if e.inf() || o.inf() {
		if e.all && o.neginf {
			ret = append(ret, o)
		} else if e.all && o.posinf {
			ret = append(ret, o)
		} else if e.posinf && o.posinf {
			ret = append(ret, &Element{first: largestOf(e.first, o.first), posinf: true})
		} else if e.posinf && o.neginf {
			ret = append(ret, &Element{first: e.first, last: o.first})
		} else if e.neginf && o.neginf {
			ret = append(ret, &Element{first: smallestOf(e.first, o.first), neginf: true})
		} else if e.neginf && o.posinf {
			ret = append(ret, &Element{first: o.first, last: e.last})
		} else if e.posinf && !o.inf() {
			ret = append(ret, &Element{first: largestOf(e.first, o.first), last: o.last})
		} else if e.neginf && !o.inf() {
			ret = append(ret, &Element{first: o.first, last: smallestOf(e.first, o.last)})
		} else if !e.inf() && o.posinf {
			ret = append(ret, &Element{first: largestOf(e.first, o.first), last: e.last})
		} else if !e.inf() && o.neginf {
			ret = append(ret, &Element{first: e.first, last: smallestOf(e.last, o.first)})
		}
	} else {
		ret = append(ret, &Element{first: largestOf(e.first, o.first), last: smallestOf(e.last, o.last)})
	}

	return ret
}

// Returns the minimum range value of two ranges. Returns two
// integers: 0 or neginf if negative infinitive and the minInt value of a
// and b unless negative infinitive.
func minInt(e, o *Element) (int, int) {
	if e.all || o.all || e.neginf || o.neginf || e.posinf || o.posinf {
		if e.neginf || o.neginf || e.all || o.all {
			return -1, 0
		} else if e.posinf && o.posinf {
			if e.first < o.first {
				return 0, e.first
			}
			return 0, o.first
		} else if e.posinf {
			return 0, o.first
		}
		return 0, e.first
	}

	if e.first < o.first {
		return 0, e.first
	}
	return 0, o.first
}

// Returns the upper limit of two ranges. The function returns two
// integers. The first is set to posinf if the maximum point is
// positive infinite. The second is set to the maxumum point unless the
// first is set to posinf.
func maxInt(e, o *Element) (int, int) {
	if e.all || o.all || e.neginf || o.neginf || e.posinf || o.posinf {
		if e.posinf || o.posinf || e.all || o.all {
			return -1, 0
		} else if e.neginf && o.neginf {
			if e.first > o.first {
				return 0, e.first
			}
			return 0, o.first
		} else if e.neginf {
			return 0, o.last
		}
		return 0, e.last
	}

	if e.last > o.last {
		return 0, e.last
	}

	return 0, o.last
}

func smallestOf(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func largestOf(a, b int) int {
	if a > b {
		return a
	}
	return b
}
