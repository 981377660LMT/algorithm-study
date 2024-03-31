// 弗罗伊德探环法(floyd探环法、floyd环检测)
// https://ferin-15.github.io/program_contest_library/library/math/floyd_cycle_find.cpp.html

package main

import "fmt"

func main() {
	// 寻找伪随机数的周期
	// 線形合同法:线性同余方法是一个产生伪随机数的方法
	s0 := 1611516670
	next := func(s int) int {
		m := 1 << 40
		return (s + (s >> 20) + 12345) & (m - 1)
	}
	start, period, ok := FloydCycleFind(s0, next, -1)
	fmt.Println(start, period, ok)
}

type C = int

// 给定一个首项为 s0 , s[i]=next(s[i-1]) (i>=1) 的序列，求环的起点和周期(长度).
//  返回值 start 为环的起点，period 为环的长度, hasCycle 为是否存在环.
//  即 s[i] = s[i+period] (i>=start).
//  !O(start+period) 时间复杂度.
func FloydCycleFind(s0 C, next func(C) C, null C) (start, period int, hasCycle bool) {
	slow := s0
	if slow == null {
		return
	}
	fast := next(slow)
	if fast == null {
		return
	}
	p1, p2 := 0, 1

	for slow != fast {
		n1 := next(fast)
		if n1 == null {
			return
		}
		n2 := next(n1)
		if n2 == null {
			return
		}
		fast = n2
		p2 += 2
		slow = next(slow)
		p1++
	}

	fast = s0
	for i := 0; i < p2-p1; i++ {
		fast = next(fast)
	}
	slow = s0
	for slow != fast {
		slow = next(slow)
		fast = next(fast)
		start++
	}
	period = 1
	for fast = next(slow); slow != fast; period++ {
		fast = next(fast)
	}
	hasCycle = true
	return
}
