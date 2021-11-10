package intset

import (
	"fmt"
	"strings"
	"testing"
)

func TestRangeJoinNegConnected(t *testing.T) {
	a := Range(-20, -10)
	b := Range(-10, 0)
	e := "-20:0"
	if fmt.Sprintf("%s", a.join(b)) != e {
		t.Fatalf("joining range: %q and %q gave %q, expected %q", a, b, a.join(b), e)
	}
}

func TestRangeJoinPosConnected(t *testing.T) {
	a := Range(10, 20)
	b := Range(20, 30)
	e := "10:30"
	if fmt.Sprintf("%s", a.join(b)) != e {
		t.Fatalf("joining range: %q and %q gave %q, expected %q", a, b, a.join(b), e)
	}
}

func TestRangeJoinCrostConnected(t *testing.T) {
	a := Range(-10, 0)
	b := Range(0, 10)
	e := "-10:10"
	if fmt.Sprintf("%s", a.join(b)) != e {
		t.Fatalf("joining range: %q and %q gave %q, expected %q", a, b, a.join(b), e)
	}
}

func TestRangeJoinNegInf(t *testing.T) {
	a := NegInf(-10)
	b := Range(-9, 1)
	e := "-∞:1"
	if fmt.Sprintf("%s", a.join(b)) != e {
		t.Fatalf("joining range: %q and %q gave %q, expected %q", a, b, a.join(b), e)
	}
}

func TestRangeJoinPosInf(t *testing.T) {
	a := PosInf(10)
	b := Range(-1, 9)
	e := "-1:∞"
	if fmt.Sprintf("%s", a.join(b)) != e {
		t.Fatalf("joining range: %q and %q gave %q, expected %q", a, b, a.join(b), e)
	}
}

func TestRangeRemovePosInf(t *testing.T) {
	a := NegInf(10)
	b := PosInf(5)
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := "-∞:4"
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

func TestRangeRemoveNegInf(t *testing.T) {
	a := PosInf(5)
	b := NegInf(10)
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := "11:∞"
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

func TestRangeRemovePosInfN(t *testing.T) {
	a := PosInf(5)
	b := Range(10, 30)
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := "31:∞"
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

func TestRangeRemoveNegInfN(t *testing.T) {
	a := NegInf(10)
	b := Range(5, 25)
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := "-∞:4"
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

// range_test.go:126: remove range: "-4:4" from "-5:5" gave "-5, -3:5", expected "-5, 5"
func TestRangeRemoveOverlap1(t *testing.T) {
	a := Range(-5, 5)
	b := Range(-4, 4)
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := "-5, 5"
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

func TestRangeRemoveOverlap2(t *testing.T) {
	a := Range(-5, 5)
	b := Range(-5, 5)
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := ""
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

func TestRangeRemoveOverlap3(t *testing.T) {
	a := Range(-5, 5)
	b := Range(-2, 8)
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := "-5:-3"
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

func TestRangeRemoveOverlap4(t *testing.T) {
	a := Range(-5, 5)
	b := Range(-10, 0)
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := "1:5"
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

func TestRangeRemoveOverlap5(t *testing.T) {
	a := All()
	b := NegInf(-10)
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := "-9:∞"
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

func TestRangeRemoveOverlap6(t *testing.T) {
	a := All()
	b := PosInf(-10)
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := "-∞:-11"
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

func TestRangeRemoveOverlap7(t *testing.T) {
	a := All()
	b := NegInf(10)
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := "11:∞"
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

func TestRangeRemoveOverlap8(t *testing.T) {
	a := PosInf(-50)
	b := All()
	r := a.remove(b)
	var s []string
	for _, re := range r {
		s = append(s, fmt.Sprintf("%s", re))
	}
	got := strings.Join(s, ", ")
	e := ""
	if got != e {
		t.Fatalf("remove range: %q from %q gave %q, expected %q", b, a, got, e)
	}
}

//

func TestIsSuper1(t *testing.T) {
	a := Range(-10, 10)
	b := Range(-9, 9)
	expect := true
	if a.isSuper(b) != expect {
		t.Fatalf("test: %q isSuper of %q returned %v, expected %v", a, b, a.isSuper(b), expect)
	}
}

func TestIsSuper2(t *testing.T) {
	a := Range(-9, 9)
	b := Range(-10, 10)
	expect := false
	if a.isSuper(b) != expect {
		t.Fatalf("test: %q isSuper of %q returned %v, expected %v", a, b, a.isSuper(b), expect)
	}
}

func TestIsSuper3(t *testing.T) {
	a := Range(-9, 9)
	b := PosInf(-10)
	expect := false
	if a.isSuper(b) != expect {
		t.Fatalf("test: %q isSuper of %q returned %v, expected %v", a, b, a.isSuper(b), expect)
	}
}

func TestIsSuper4(t *testing.T) {
	a := PosInf(-9)
	b := Range(-8, 22)
	expect := true
	if a.isSuper(b) != expect {
		t.Fatalf("test: %q isSuper of %q returned %v, expected %v", a, b, a.isSuper(b), expect)
	}
}

func TestIsSuper5(t *testing.T) {
	a := NegInf(10)
	b := PosInf(5)
	e := false
	if a.isSuper(b) != e {
		t.Fatalf("test: %q isSuper of %q returned %v, expected %v", a, b, a.isSuper(b), e)
	}
}

func TestIsAdjacent1a(t *testing.T) {
	a := All()
	b := PosInf(5)
	e := true
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent1b(t *testing.T) {
	a := PosInf(5)
	b := All()
	e := true
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent2a(t *testing.T) {
	a := NegInf(4)
	b := Range(5, 10)
	e := true
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent2b(t *testing.T) {
	a := Range(5, 10)
	b := NegInf(4)
	e := true
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent2a1(t *testing.T) {
	a := NegInf(5)
	b := Range(5, 10)
	e := false
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent3a(t *testing.T) {
	a := PosInf(11)
	b := Range(5, 10)
	e := true
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent3b(t *testing.T) {
	a := Range(5, 39)
	b := PosInf(40)
	e := true
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent4a(t *testing.T) {
	a := NegInf(-5)
	b := PosInf(-4)
	e := true
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent4b(t *testing.T) {
	b := PosInf(-4)
	a := NegInf(-5)
	e := true
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent5a(t *testing.T) {
	a := Range(-5, 10)
	b := Range(11, 20)
	e := true
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent5b(t *testing.T) {
	a := Range(11, 20)
	b := Range(-5, 10)
	e := true
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent5a1(t *testing.T) {
	a := Range(-50, 10)
	b := Range(10, 20)
	e := false
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent5b1(t *testing.T) {
	b := Range(-50, 10)
	a := Range(10, 20)
	e := false
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsAdjacent6b1(t *testing.T) {
	b := Range(59, 499)
	a := PosInf(500)
	e := true
	if a.isAdjacent(b) != e {
		t.Fatalf("test: %q isAdjacent of %q returned %v, expected %v", a, b, a.isAdjacent(b), e)
	}
}

func TestIsOverlapping1a(t *testing.T) {
	a := All()
	b := Range(10, 20)
	e := true
	if a.isOverlapping(b) != e {
		t.Fatalf("test: %q isOverlapping of %q returned %v, expected %v", a, b, a.isOverlapping(b), e)
	}
}

func TestIsOverlapping1b(t *testing.T) {
	a := Range(10, 20)
	b := All()
	e := true
	if a.isOverlapping(b) != e {
		t.Fatalf("test: %q isOverlapping of %q returned %v, expected %v", a, b, a.isOverlapping(b), e)
	}
}

func TestIsOverlapping2a(t *testing.T) {
	a := Range(-10, -5)
	b := NegInf(5)
	e := true
	if a.isOverlapping(b) != e {
		t.Fatalf("test: %q isOverlapping of %q returned %v, expected %v", a, b, a.isOverlapping(b), e)
	}
}

func TestIsOverlapping2b(t *testing.T) {
	a := NegInf(5)
	b := Range(-10, -5)
	e := true
	if a.isOverlapping(b) != e {
		t.Fatalf("test: %q isOverlapping of %q returned %v, expected %v", a, b, a.isOverlapping(b), e)
	}
}

func TestIsOverlapping3a(t *testing.T) {
	a := Range(-10, 5)
	b := PosInf(5)
	e := true
	if a.isOverlapping(b) != e {
		t.Fatalf("test: %q isOverlapping of %q returned %v, expected %v", a, b, a.isOverlapping(b), e)
	}
}

func TestIsOverlapping3b(t *testing.T) {
	a := PosInf(5)
	b := Range(5, 10)
	e := true
	if a.isOverlapping(b) != e {
		t.Fatalf("test: %q isOverlapping of %q returned %v, expected %v", a, b, a.isOverlapping(b), e)
	}
}

func TestIsOverlapping4a(t *testing.T) {
	a := Range(5, 10)
	b := Range(2, 8)
	e := true
	if a.isOverlapping(b) != e {
		t.Fatalf("test: %q isOverlapping of %q returned %v, expected %v", a, b, a.isOverlapping(b), e)
	}
}

func TestIsOverlapping4b(t *testing.T) {
	a := Range(2, 8)
	b := Range(5, 10)
	e := true
	if a.isOverlapping(b) != e {
		t.Fatalf("test: %q isOverlapping of %q returned %v, expected %v", a, b, a.isOverlapping(b), e)
	}
}

func TestIsOverlapping5a(t *testing.T) {
	a := Range(5, 10)
	b := Range(10, 18)
	e := true
	if a.isOverlapping(b) != e {
		t.Fatalf("test: %q isOverlapping of %q returned %v, expected %v", a, b, a.isOverlapping(b), e)
	}
}

func TestIsOverlapping5b(t *testing.T) {
	a := Range(10, 18)
	b := Range(5, 10)
	e := true
	if a.isOverlapping(b) != e {
		t.Fatalf("test: %q isOverlapping of %q returned %v, expected %v", a, b, a.isOverlapping(b), e)
	}
}
