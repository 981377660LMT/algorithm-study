// https://www.luogu.com.cn/problem/P3391

package main

import (
	"bufio"
	"fmt"
	"os"
)

func 文艺平衡树() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i + 1
	}

	t := NewFHQTreap(nums)

	for i := 0; i < m; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		t.Reverse(l-1, r)
	}

	// nlogn 输出
	// for i := 0; i < n; i++ {
	// 	fmt.Fprint(out, t.At(i), " ")
	// }

	// !O(n) 中序遍历输出
	res := t.InOrder()
	for i := 0; i < n; i++ {
		fmt.Fprint(out, res[i], " ")
	}
}
