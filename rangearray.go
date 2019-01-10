package rangearray

// Package rangearray implements indexing and searching in a semi-dense
// array.  In particular, a semi-dense array has runs of consecutive
// entries with potentially large (and irregular) separation between
// those runs.
//
// Note that the data structure is optimized for insertions only at the
// end; insertions before the end can be very expensive.  A search over
// n runs takes O(log(n)) time.

// As motivation, the original application for this package was to
// manage observation data for GPS satellites.  A satellite "pass" is
// typically several hours of consecutive data, but any GPS satellite
// only has two passes (over one ground station) per day.  The time of
// each pass shifts slightly from day to day, and there are sometimes
// glitches that cause missing data, so there is no clever way to know
// the number of valid observations before time T in constant time
// (without a laughably inefficient or complex data structure).
//
// A lot of interesting data series are too large to fit in a reasonable
// desktop's memory, but interesting blocks of observables will -- and
// time is the common factor to relate them.  This suggests using a
// structure-of-arrays representation with a very dense index of times.
// (Conveniently, the observables are collected at fixed intervals that
// are multiples of one second, so the times can be easily converted to
// uint32 values that cover a 136-year span.)

// A rangearray is stored using a run-length encoding (RLE) format.
//
// Each RLE entry has three elements: value, index, and count, which
// are the starting value of this run, the number of entries in the
// rangearray before this run, and the number of consecutive entries
// in this run.
//
// Each rangearray is a slice of RLE entries.

import (
	"sort"
)

// Uint32Run is an RLE entry in a Uint32 rangearray.
type Uint32Run struct {
	// Value is the starting value of this run.
	Value uint32

	// Index is the number of elements before this run.
	Index uint32

	// Count is the number of consecutive elements inside this run.
	Count uint32
}

// Uint32 is a semi-dense array of Uint32 values.  The zero value is an
// empty rangearray.
type Uint32 struct {
	S []Uint32Run
}

// Min returns the minimum value in r.  Panics if r is empty.
func (r Uint32) Min() uint32 {
	return r.S[0].Value
}

// Max returns the maximum value in r.  Panics if r is empty.
func (r Uint32) Max() uint32 {
	n := len(r.S) - 1
	return r.S[n].Value + r.S[n].Count - 1
}

// Len returns the number of elements in r.
func (r Uint32) Len() uint32 {
	if r.S == nil {
		return 0
	}

	n := len(r.S) - 1
	return r.S[n].Index + r.S[n].Count
}

// IndexOf returns the number of elements in r that are less than x.
func (r Uint32) IndexOf(x uint32) uint32 {
	// Common case: x <= r.Max().
	i := r.LowerBound(x)
	if i < len(r.S) {
		if x <= r.S[i].Value {
			return r.S[i].Index
		}
		return x - r.S[i].Value + r.S[i].Index
	}

	// Otherwise, r is empty or x > r.Max().
	return r.Len()
}

// LowerBound returns the index of the run in r that contains x.  If no
// run contains x, LowerBound returns the index of the run that starts
// after x.  If x is after r.Max(), returns len(r.S).
func (r Uint32) LowerBound(x uint32) int {
	if r.S == nil {
		return 0
	}

	return sort.Search(len(r.S), func(i int) bool {
		return x < r.S[i].Value+r.S[i].Count
	})
}

// Push adds x to r.
func (r *Uint32) Push(x uint32) {
	// Is this the first entry?
	if r.S == nil {
		r.S = append(r.S, Uint32Run{
			Value: x,
			Index: 0,
			Count: 1,
		})
		return
	}

	// Can we append to the last entry?
	n := len(r.S) - 1
	if r.S[n].Value+r.S[n].Count == x {
		r.S[n].Count++
		return
	}

	// Is it past the last entry?
	if r.S[n].Value+r.S[n].Count < x {
		r.S = append(r.S, Uint32Run{
			Value: x,
			Index: r.S[n].Index + r.S[n].Count,
			Count: 1,
		})
		return
	}

	// Find the insertion point.
	n = r.LowerBound(x)
	if x >= r.S[n].Value {
		// either x is within r.S[n] and we silently ignore the dupe...
		// or x is after r.S[n] and LowerBound() had a bug
		return
	}

	// Is x just after r.S[n-1]?
	afterNm1 := n > 0 && x == r.S[n-1].Value+r.S[n-1].Count

	// Is x just before r.S[n]?
	if x+1 == r.S[n].Value {
		if afterNm1 {
			// Merge r.S[n] into r.S[n-1] and shrink the rest.
			r.S[n-1].Count += r.S[n].Count + 1
			copy(r.S[n:], r.S[n+1:])
			r.S = r.S[:len(r.S)-1]
			n--
		} else {
			r.S[n].Value--
			r.S[n].Count++
		}
	} else if afterNm1 {
		r.S[n-1].Count++
		n--
	} else {
		l := len(r.S)
		r.S = append(r.S, r.S[l-1])
		copy(r.S[n+1:l], r.S[n:l-1])
		r.S[n] = Uint32Run{
			Value: x,
			Index: r.S[n+1].Index,
			Count: 1,
		}
	}

	for n+1 < len(r.S) {
		n++
		r.S[n].Index++
	}
}
