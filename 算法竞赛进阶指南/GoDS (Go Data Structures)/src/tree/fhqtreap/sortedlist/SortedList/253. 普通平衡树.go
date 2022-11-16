package main

import (
	"bufio"
	"fmt"
	"os"
)

// 您需要写一种数据结构（可参考题目标题），来维护一些数，其中需要提供以下操作：
// 插入数值 x。
// 删除数值 x(若有多个相同的数，应只删除一个)。
// 查询数值 x 的排名(若有多个相同的数，应输出最小的排名)。
// 查询排名为 x 的数值。
// 求数值 x 的前驱(前驱定义为小于 x 的最大的数)。
// 求数值 x 的后继(后继定义为大于 x 的最小的数)。
func demo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	sl := NewSortedList(func(a, b Value) int {
		return a.(int) - b.(int)
	}, q)

	for i := 0; i < q; i++ {
		var op, x int
		fmt.Fscan(in, &op, &x)
		switch op {
		case 1:
			sl.Add(x)
		case 2:
			sl.Discard(x)
		case 3:
			res := sl.BisectLeft(x) + 1
			fmt.Fprintln(out, res)
		case 4:
			x--
			res := sl.At(x)
			fmt.Fprintln(out, res)
		case 5:
			// 前驱
			pos := sl.BisectLeft(x) - 1
			res := sl.At(pos)
			fmt.Fprintln(out, res)
		case 6:
			// 后继
			pos := sl.BisectRight(x)
			res := sl.At(pos)
			fmt.Fprintln(out, res)
		}
	}
}
