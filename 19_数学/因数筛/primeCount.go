// https://rsk0315.hatenablog.com/entry/2021/05/18/015511

// !O(n^(3/4)/logn)
// 不超过n的素数个数 n<=1e12

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	fmt.Fprintln(out, CountPrime(n))
}

func CountPrime(n int) int {
	if n <= 1 {
		return 0
	}
	if n == 2 {
		return 1
	}

	v := isqrt(n)
	s := (v + 1) >> 1
	smalls, roughs, larges := make([]int, s), make([]int, s), make([]int, s)
	for i := 1; i < s; i++ {
		smalls[i] = i
	}
	for i := 0; i < s; i++ {
		roughs[i] = i<<1 | 1
		larges[i] = (n/(i<<1|1) - 1) >> 1
	}
	skip := make([]bool, v+1)

	divide := func(n, d int) int { return n / d }
	half := func(n int) int { return (n - 1) >> 1 }

	pc := 0
	for p := 3; p <= v; p += 2 {
		if skip[p] {
			continue
		}
		q := p * p
		if q*q > n {
			break
		}
		skip[p] = true
		for i := q; i <= v; i += 2 * p {
			skip[i] = true
		}
		ns := 0
		for k := 0; k < s; k++ {
			i := roughs[k]
			if skip[i] {
				continue
			}
			d := i * p
			if d <= v {
				larges[ns] = larges[k] - larges[smalls[d>>1]-pc] + pc
			} else {
				larges[ns] = larges[k] - smalls[half(divide(n, d))] + pc
			}
			roughs[ns] = i
			ns++
		}
		s = ns
		for i, j := half(v), ((v/p)-1)|1; j >= p; j -= 2 {
			c := smalls[j>>1] - pc
			for e := (j * p) >> 1; i >= e; i-- {
				smalls[i] -= c
			}
		}
		pc++
	}

	larges[0] += (s + 2*(pc-1)) * (s - 1) >> 1
	for k := 1; k < s; k++ {
		larges[0] -= larges[k]
	}
	for l := 1; l < s; l++ {
		q := roughs[l]
		M := n / q
		e := smalls[half(divide(M, q))] - pc
		if e < l+1 {
			break
		}
		t := 0
		for k := l + 1; k <= e; k++ {
			t += smalls[half(divide(M, roughs[k]))]
		}
		larges[0] += t - (e-l)*(pc+l-1)
	}

	return larges[0] + 1
}

func isqrt(x int) int {
	x0 := x >> 1
	if x0 == 0 {
		return x
	}
	x1 := (x0 + x/x0) >> 1
	for x1 < x0 {
		x0 = x1
		x1 = (x0 + x/x0) >> 1
	}
	return x0
}
