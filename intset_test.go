package intset

import (
	"fmt"
	"testing"
)

func TestIntSetCompose1(t *testing.T) {
	a := New()
	e := "{∅}"
	if fmt.Sprintf("%s", a) != e {
		t.Fatalf("compose failed: got %s, expected %s", a, e)
	}
}

func TestIntSetCompose2(t *testing.T) {
	a := New(All())
	e := "{-∞:∞}"
	if fmt.Sprintf("%s", a) != e {
		t.Fatalf("compose failed: got %s, expected %s", a, e)
	}
}

func TestIntSetCompose3(t *testing.T) {
	a := New(PosInf(-4))
	e := "{-4:∞}"
	if fmt.Sprintf("%s", a) != e {
		t.Fatalf("compose failed: got %s, expected %s", a, e)
	}
}

func TestIntSetCompose4(t *testing.T) {
	a := New(NegInf(4))
	e := "{-∞:4}"
	if fmt.Sprintf("%s", a) != e {
		t.Fatalf("compose failed: got %s, expected %s", a, e)
	}
}

func TestIntSetCompose5(t *testing.T) {
	a := New(Range(-400, -200), Range(-150, -34), Range(400, 420), Range(50, 64), Range(90, 100), PosInf(500))
	e := "{-400:-200, -150:-34, 50:64, 90:100, 400:420, 500:∞}"
	if fmt.Sprintf("%s", a) != e {
		t.Fatalf("compose failed: got %s, expected %s", a, e)
	}
}

func TestIntSetCompose6(t *testing.T) {
	a := New(Range(-400, -200), Range(-199, -34), Range(400, 420), Range(50, 399), Range(49, 101), PosInf(500), NegInf(-5000))
	e := "{-∞:-5000, -400:-34, 49:420, 500:∞}"
	if fmt.Sprintf("%s", a) != e {
		t.Fatalf("compose failed: got %s, expected %s", a, e)
	}
}

func TestIntSetRemove1(t *testing.T) {
	a := New(All())
	a.RemoveElements(PosInf(5000), NegInf(-5000))
	e := "{-4999:4999}"
	if fmt.Sprintf("%s", a) != e {
		t.Fatalf("remove failed: got %s, expected %s", a, e)
	}
}

func TestIntSetRemove2(t *testing.T) {
	a := New(All())
	a.RemoveElements(PosInf(500), NegInf(-500))
	a.RemoveInts(12, 14, 16)
	e := "{-499:11, 13, 15, 17:499}"
	if fmt.Sprintf("%s", a) != e {
		t.Fatalf("remove failed: got %s, expected %s", a, e)
	}
}

func TestComplement1(t *testing.T) {
	a := New(All())
	e := "{∅}"
	if fmt.Sprintf("%s", a.Complement()) != e {
		t.Fatalf("complement failed: %s%c, got %s, expected %s", a, 0x2201, a.Complement(), e)
	}
}

func TestComplement2(t *testing.T) {
	a := New()
	e := "{-∞:∞}"
	if fmt.Sprintf("%s", a.Complement()) != e {
		t.Fatalf("complement failed: %s%c, got %s, expected %s", a, 0x2201, a.Complement(), e)
	}
}

func TestComplement3(t *testing.T) {
	a := New()
	e := "{∅}"
	if fmt.Sprintf("%s", a.Complement().Complement()) != e {
		t.Fatalf("complement failed: %s%c%c, got %s, expected %s", a, 0x2201, 0x2201, a.Complement().Complement(), e)
	}
}

func TestComplement4(t *testing.T) {
	a := New(Range(-400, -200), Range(-199, -34), Range(400, 420), Range(50, 399), Range(49, 101), PosInf(500), NegInf(-5000))
	e := "{-4999:-401, -33:48, 421:499}"
	if fmt.Sprintf("%s", a.Complement()) != e {
		t.Fatalf("complement failed:  %s%c, got %s, expected %s", a, 0x2201, a.Complement(), e)
	}
}

func TestUnion1(t *testing.T) {
	a := New(Range(-400, -200), Range(-199, -34), Range(400, 420), Range(50, 399), Range(49, 101), PosInf(500), NegInf(-5000))
	b := New(Range(-400, -200), Range(-199, -34), Range(400, 420), Range(50, 399), Range(49, 101), PosInf(500), NegInf(-5000))
	e := "{-∞:-5000, -400:-34, 49:420, 500:∞}"
	if fmt.Sprintf("%s", a.Union(b)) != e {
		t.Fatalf("union failed: %s %c %s, got %s, expected %s", a, 0x222a, b, a.Union(b), e)
	}
}

func TestUnion2(t *testing.T) {
	a := New(Range(-400, -200), Range(-199, -34), Range(400, 420), Range(50, 399), Range(49, 101), PosInf(500), NegInf(-5000))
	b := a.Complement()
	e := "{-∞:∞}"
	if fmt.Sprintf("%s", a.Union(b)) != e {
		t.Fatalf("union failed: %s %c %s, got %s, expected %s", a, 0x222a, b, a.Union(b), e)
	}
}

func TestIntersect1(t *testing.T) {
	a := New(NegInf(-5), PosInf(5))
	b := New(Range(-50, 50))
	e := "{-50:-5, 5:50}"
	if fmt.Sprintf("%s", a.Intersect(b)) != e {
		t.Fatalf("intersection failed: %s %c %s, got %s, expected %s", a, 0x2229, b, a.Intersect(b), e)
	}
}

func TestDifference1(t *testing.T) {
	a := New(Range(-90, 5), Range(90, 100))
	b := New(Range(-50, 50))
	e := "{-90:-51, 90:100}"
	if fmt.Sprintf("%s", a.Difference(b)) != e {
		t.Fatalf("difference failed: %s - %s, got %s, expected %s", a, b, a.Difference(b), e)
	}
}

func TestCardinality1a(t *testing.T) {
	maxuint := ^uint(0)
	maxint := int(maxuint >> 1)
	minint := -maxint - 1
	a := New(Range(minint, maxint))
	c, inf := a.Cardinality()
	if !inf {
		t.Fatalf("cardinality failed: %s, got %d, expected %c", a, c, 0x221e)
	}
}

func TestCardinality1b(t *testing.T) {
	maxuint := ^uint(0)
	maxint := int(maxuint >> 1)
	minint := -maxint
	a := New(Range(minint, maxint))
	var e = ^uint(0)
	c, inf := a.Cardinality()
	if inf {
		t.Fatalf("cardinality failed: %s, got %c, expected %d", a, 0x221e, e)
	} else if c != e {
		t.Fatalf("cardinality failed: %s, got %d, expected %d", a, c, e)
	}
}

func TestCardinality2(t *testing.T) {
	a := New(Range(-1, 1))
	var e uint = 3
	c, inf := a.Cardinality()
	if inf {
		t.Fatalf("cardinality failed: %s, got %c, expected %d", a, 0x221e, e)
	} else if c != e {
		t.Fatalf("cardinality failed: %s, got %d, expected %d", a, c, e)
	}
}

func TestCardinality3(t *testing.T) {
	a := New(Range(1, 5))
	var e uint = 5
	c, inf := a.Cardinality()
	if inf {
		t.Fatalf("cardinality failed: %s, got %c, expected %d", a, 0x221e, e)
	} else if c != e {
		t.Fatalf("cardinality failed: %s, got %d, expected %d", a, c, e)
	}
}

func TestCardinality4(t *testing.T) {
	a := New(Range(-5, -1))
	var e uint = 5
	c, inf := a.Cardinality()
	if inf {
		t.Fatalf("cardinality failed: %s, got %c, expected %d", a, 0x221e, e)
	} else if c != e {
		t.Fatalf("cardinality failed: %s, got %d, expected %d", a, c, e)
	}
}

func TestXor(t *testing.T) {
	a := New(Range(-10, -5), Range(5, 10), PosInf(25))
	b := New(Range(-8, -3), Range(2, 6))

	e := "{-10:-9, -4:-3, 2:4, 7:10, 25:∞}"
	if fmt.Sprintf("%s", a.Xor(b)) != e {
		t.Fatalf("xor failed: %s %c %s, got %s, expected %s", a, 0x22bb, b, a.Xor(b), e)
	}
}
