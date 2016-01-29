// Package ddmin implements the ddmin test minimization algorithm
/*

Simplifying and Isolating Failure-Inducing Input
Andreas Zeller (2002)

    https://www.st.cs.uni-saarland.de/papers/tse2002/tse2002.pdf

*/
package ddmin

import (
	"math"
)

type Result int

const (
	// Pass indicates the test passed
	Pass Result = iota
	// Fail indicates the expected test failure was produced
	Fail
	// Unresolved indicates the test failed for a different reason
	Unresolved
)

// looks to minimize data so that f will fail
func Minimize(data []byte, f func(d []byte) Result) []byte {

	if f(nil) == Fail {
		// that was easy..
		return nil
	}

	if f(data) == Pass {
		panic("ddmin: function must fail on data")
	}

	return ddmin(data, f, 2)
}

func ddmin(data []byte, f func(d []byte) Result, granularity int) []byte {

mainloop:
	for len(data) >= 2 {

		subsets := makeSubsets(data, granularity)

		for _, subset := range subsets {
			if f(subset) == Fail {
				// fake tail recursion
				data = subset
				granularity = 2
				continue mainloop
			}
		}

		b := make([]byte, len(data))
		for i := range subsets {
			complement := makeComplement(subsets, i, b[:0])
			if f(complement) == Fail {
				granularity--
				if granularity < 2 {
					granularity = 2
				}
				// fake tail recursion
				data = complement
				continue mainloop
			}
		}

		if granularity == len(data) {
			return data
		}

		granularity *= 2

		if granularity > len(data) {
			granularity = len(data)
		}
	}

	return data
}

func makeSubsets(data []byte, granularity int) [][]byte {

	var subsets [][]byte

	// Use a variation of https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm to generate equal(ish) sized subsets
	// TODO(dgryski): maybe switch to integer algorithm instead of floating point

	fsize := float64(len(data)) / float64(granularity)
	isize, frac := int(fsize), fsize-math.Trunc(fsize)

	var ferr float64
	for i := 0; i < granularity-1; i++ {
		ferr += frac
		size := isize
		if ferr > 0.5 {
			size++
			ferr -= 1.0
		}
		subsets = append(subsets, data[:size])
		data = data[size:]
	}
	subsets = append(subsets, data)

	return subsets
}

func makeComplement(subsets [][]byte, n int, b []byte) []byte {
	for i, s := range subsets {
		if i == n {
			continue
		}
		b = append(b, s...)
	}
	return b
}
