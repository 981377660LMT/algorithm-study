package fhqtreap

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	nums := make([]interface{}, n)
	for i := range nums {
		nums[i] = i + 1
	}

	t := NewFHQTreap(nums, func(a, b interface{}) int {
		return a.(int) - b.(int)
	})

	for i := 0; i < m; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		t.Reverse(l-1, r-1)
	}

	fmt.Fprintln(out, t.InOrder()...)
}
