package rangearray

import (
	"testing"
)

func testEmptyMin(t *testing.T, r Uint32) {
	defer func() {
		if recover() == nil {
			t.Errorf("Expected Uint32{}.Min() to panic, but it didn't")
		}
	}()

	_ = r.Min()
}

func testEmptyMax(t *testing.T, r Uint32) {
	defer func() {
		if recover() == nil {
			t.Errorf("Expected Uint32{}.Max() to panic, but it didn't")
		}
	}()

	_ = r.Max()
}

func TestEmptyUint32(t *testing.T) {
	r := Uint32{}

	testEmptyMin(t, r)
	testEmptyMax(t, r)
	if x := r.Len(); x != 0 {
		t.Errorf("Expected Uint32{}.Len() == 0, got %d", x)
	}
	if x := r.IndexOf(0); x != 0 {
		t.Errorf("Expected Uint32{}.IndexOf(0) == 0, got %d", x)
	}
	if x := r.LowerBound(0); x != 0 {
		t.Errorf("Expected Uint32{}.LowerBound(0) == 0, got %d", x)
	}
}

type indexOfUint32 struct {
	value, index uint32
}

func testIndexOf(t *testing.T, r Uint32, v []indexOfUint32) {
	for _, s := range v {
		if x := r.IndexOf(s.value); x != s.index {
			t.Errorf("Expected r.IndexOf(%d) == %d, got %d", s.value, s.index, x)
		}
	}
}

func TestSearchUint32(t *testing.T) {
	r := &Uint32{}

	for i := 100; i < 200; i++ {
		r.Push(uint32(i))
	}
	for i := 350; i < 450; i++ {
		r.Push(uint32(i))
	}

	if len(r.S) != 2 {
		t.Errorf("Expected len(r.S) == 2, got %d", len(r.S))
	}

	if x := r.LowerBound(50); x != 0 {
		t.Errorf("Expected r.LowerBound(50) == 0, got %d", x)
	}

	testIndexOf(t, *r, []indexOfUint32{
		{275, 100},
		{400, 150},
		{500, 200},
	})

	r.Push(75)
	r.Push(200)
	r.Push(202)
	r.Push(201)
	r.Push(200)
	r.Push(349)

	testIndexOf(t, *r, []indexOfUint32{
		{75, 0},
		{100, 1},
		{200, 101},
		{202, 103},
		{203, 104},
		{349, 104},
		{350, 105},
		{450, 205},
		{500, 205},
	})
}
