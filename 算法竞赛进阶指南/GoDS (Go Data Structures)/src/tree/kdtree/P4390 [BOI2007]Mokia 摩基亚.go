// 矩形区域面积<=2e5*2e5

package main

import (
	"bufio"
	"fmt"
	"os"
)

func demo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	tree := NewKdTree()
	var kind, w, x1, y1, val, x2, y2 int
	fmt.Fscan(in, &kind, &w)
	for {
		fmt.Fscan(in, &kind)
		if kind == 3 {
			break
		}
		if kind == 1 {
			// 向方格 (x,y) 中添加 a 个用户。a 是正整数。
			fmt.Fscan(in, &x1, &y1, &val)
			tree.Add([2]int{x1, y1}, val)
		} else {
			// 查询 所规定的矩形中的用户数量。
			fmt.Fscan(in, &x1, &y1, &x2, &y2)
			fmt.Fprintln(out, tree.Query(x1, y1, x2, y2))
		}
	}
}
