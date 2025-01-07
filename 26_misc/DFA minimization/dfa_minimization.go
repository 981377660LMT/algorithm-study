// Package mindfa implements DFA minimization using Hopcroft's algorithm.
// https://en.wikipedia.org/wiki/DFA_minimization
// Hopcroft's algorithm
//

package main

import (
	"fmt"
	"slices"
)

func main() {
	nState := 6
	nSymbol := 2
	finals := []int{2, 3, 4}
	transitions := [][]int{
		//  0  1
		0: {1, 2}, // a
		1: {0, 3}, // b
		2: {4, 5}, // c
		3: {4, 5}, // d
		4: {4, 5}, // e
		5: {5, 5}, // f
	}
	transitionFunc := func(state, symbol int) int { return transitions[state][symbol] }

	partitions := DfaMinimization(nState, nSymbol, finals, transitionFunc)
	fmt.Println(partitions)
}

// Minimize 接受一个 DFA 表示，进行最小化处理并返回状态分组。
// 同一分组中的状态具有相同行为。
//
// 所有参数代表一个单一的 DFA（最小化前）。
// 它包含0..nState 个状态，0..nSymbol 个输入符号，并通过 transition 函数定义转移关系。
//
// transition 是一个函数，接受状态和符号并返回目标状态。
// finals 中包含的状态是接受状态。结果分区的顺序不确定，但每个分区内的状态顺序是升序的。
//
// 它使用 Hopcroft 算法。内存复杂度为 O(nState)。
//
// finals 中的所有数字必须位于 [0, nState) 范围内，且不能出现重复值。
func DfaMinimization(nState, nSymbol int, finals []int, transition func(state, symbol int) int) [][]int {
	if len(finals) > nState {
		panic(fmt.Sprintf("len(finals) should be less than or equal to nState: len(finals) = %d, nState = %d", len(finals), nState))
	}

	whole := make([]int, nState+max(len(finals), nState-len(finals)))
	copy(whole, finals)
	slices.Sort(whole[:len(finals)])

	for i := 0; i < len(finals)-1; i++ {
		if whole[i] == whole[i+1] {
			panic(fmt.Sprintf("finals contains same two value: %d", whole[i]))
		}
	}

	cmpl(whole[len(finals):nState], whole[:len(finals)], nState)

	buf := whole[nState:]

	// partitions：初始化为两组：
	// 接受状态组：whole[:len(finals)]
	// 非接受状态组：whole[len(finals):nState]
	partitions := [][]int{whole[:len(finals)], whole[len(finals):nState]}
	// works is a set of the partition which has never tried to be split.
	works := [][]int{whole[:len(finals)], whole[len(finals):nState]} // map?

	for len(works) > 0 {
		for c := 0; c < nSymbol; c++ {
			for ip, pFrom := range partitions {
				ip1, ip2 := 0, len(buf)-1
				for _, state := range pFrom {
					if includes(works[0], transition(state, c)) {
						buf[ip1] = state
						ip1++
					} else {
						buf[ip2] = state
						ip2--
					}
				}

				if ip1 == 0 || ip2 == len(buf)-1 {
					continue
				}

				p1 := pFrom[:ip1]
				copy(p1, buf[:ip1])

				p2 := pFrom[ip1:]
				for i := range p2 {
					p2[i] = buf[len(buf)-1-i]
				}

				var split bool
				for i, w := range works {
					if &w[0] != &pFrom[0] {
						continue
					}

					// Split works[i].
					works = append(works, p2)
					works[i] = p1
					split = true
					break
				}

				if !split {
					if len(p1) < len(p2) {
						works = append(works, p1)
					} else {
						works = append(works, p2)
					}
				}
				partitions[ip] = p1
				partitions = append(partitions, p2) // Don't worry, p2 is not iterated in the current loop.
			}
		}
		// pseudo-shift
		works[0] = works[len(works)-1]
		works = works[:len(works)-1]
	}
	return partitions
}

// cmpl returns the complement set of a in (0..upper).
func cmpl(dst, a []int, upper int) {
	var n, i int
	for _, u := range a {
		for ; n < u; n++ {
			dst[i] = n
			i++
		}
		n++
	}
	for ; n < upper; n++ {
		dst[i] = n
		i++
	}
}

func includes(a []int, e int) bool {
	_, ok := slices.BinarySearch(a, e)
	return ok
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
