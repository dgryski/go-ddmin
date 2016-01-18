// Package ddmin implements the ddmin test minimization algorithm
/*

Simplifying and Isolating Failure-Inducing Input
Andreas Zeller (2002)

    https://www.st.cs.uni-saarland.de/papers/tse2002/tse2002.pdf

*/
package ddmin

import (
	"bytes"
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

	var b bytes.Buffer

	for len(data) >= 2 {

		subsets := makeSubsets(data, granularity)

		for _, subset := range subsets {
			if f(subset) == Fail {
				// recurse
				return ddmin(subset, f, 2)
			}
		}

		for i := range subsets {
			complement := makeComplement(subsets, i, &b)
			if f(complement) == Fail {
				granularity--
				if granularity < 2 {
					granularity = 2
				}
				return ddmin(complement, f, granularity)
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

	size := len(data) / granularity
	for i := 0; i < granularity; i++ {
		subsets = append(subsets, data[:size])
		data = data[size:]
	}

	return subsets
}

func makeComplement(subsets [][]byte, n int, b *bytes.Buffer) []byte {

	b.Reset()

	for i, s := range subsets {
		if i == n {
			continue
		}
		b.Write(s)
	}

	return b.Bytes()
}
