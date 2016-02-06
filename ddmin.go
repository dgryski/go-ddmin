// Package ddmin implements the ddmin test minimization algorithm
/*

Simplifying and Isolating Failure-Inducing Input
Andreas Zeller (2002)

    https://www.st.cs.uni-saarland.de/papers/tse2002/tse2002.pdf

*/
package ddmin

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

func makeSubsets(s []byte, n int) [][]byte {

	// via https://www.reddit.com/r/golang/comments/44cl7f/a_better_subslicing_algorithm/czp9r6j

	ret := make([][]byte, 0, n)

	for ; n > 0; n-- {
		i := len(s) / n
		s, ret = s[i:], append(ret, s[:i])
	}

	return ret
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
