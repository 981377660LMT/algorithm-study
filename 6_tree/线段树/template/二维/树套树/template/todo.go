package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

}

// P3332 [ZJOI2013] K大数查询
// https://www.luogu.com.cn/problem/P3332
// 初识时有n个空数组，每个数组有一个编号，编号从0到n-1.
// 1 start end c :将 c 加入编号在 [start,end) 的数组中.
// 2 start end k :查询编号在 [start,end) 的数组的所有数中第 k 大的数, k>=1.
func K大数查询() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	queries := make([][4]int, q)
	for i := 0; i < q; i++ {
		var op, start, end, c int
		fmt.Fscan(in, &op, &start, &end, &c)
		start--
		queries[i] = [4]int{op, start, end, c}
	}
}

// P2617 Dynamic Rankings
// https://www.luogu.com.cn/problem/P2617
// 动态区间第k小
