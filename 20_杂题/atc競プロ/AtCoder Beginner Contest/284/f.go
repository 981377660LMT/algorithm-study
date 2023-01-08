package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

func main() {
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var t string
	fmt.Fscan(in, &t)

	getChar1 := func(insertPos, i int) byte {
		if i < insertPos || i >= n+insertPos {
			return t[i]
		}
		return t[n+i]
	}

	getChar2 := func(insertPos, i int) byte {
		end := n + insertPos
		return t[end-i-1]
	}

	isSame := func(insertPos int) bool {
		for i := 0; i < n; i++ {
			if getChar1(insertPos, i) != getChar2(insertPos, i) {
				return false
			}
		}
		return true
	}

	fail := 0
	for i := 0; i <= n; i++ {
		ok := true
		for j := 0; j < 100; j++ {
			x := rand.Intn(n)
			if getChar1(i, x) != getChar2(i, x) {
				ok = false
				break
			}
		}

		if ok {
			if isSame(i) {
				fmt.Fprintln(out, t[:i]+t[n+i:])
				fmt.Fprintln(out, i)
				return
			} else {
				fail++
			}
		}

		if fail >= 300 {
			fmt.Fprintln(out, -1)
			return
		}
	}

	fmt.Fprintln(out, -1)

}
