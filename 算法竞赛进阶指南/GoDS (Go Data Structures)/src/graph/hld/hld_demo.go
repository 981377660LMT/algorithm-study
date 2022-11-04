package hld

// https://www.acwing.com/problem/content/2570/
// 1 root1 root2 add，表示将树从 root1 到 root2 结点最短路径上所有节点的值都加上 add。
// 2 root1 add，表示将以 root1 为根节点的子树内所有节点值都加上 add。
// 3 root1 root2，表示求树从 root1 到 root2 结点最短路径上所有节点的值之和。
// 4 root1 表示求以 root1 为根节点的子树内所有节点值之和

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

	vals := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &vals[i])
	}

	adjList := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		adjList[u] = append(adjList[u], v)
		adjList[v] = append(adjList[v], u)
	}

	hld := HeavyLightDecomposition(n, 0, adjList, vals)
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op, root1, root2, add int
		fmt.Fscan(in, &op)
		switch op {
		case 1:
			fmt.Fscan(in, &root1, &root2, &add)
			root1, root2 = root1-1, root2-1
			hld.UpdatePath(root1, root2, add)
		case 2:
			fmt.Fscan(in, &root1, &add)
			root1--
			hld.UpdateSubtree(root1, add)
		case 3:
			fmt.Fscan(in, &root1, &root2)
			root1, root2 = root1-1, root2-1
			fmt.Fprintln(out, hld.QueryPath(root1, root2))
		case 4:
			fmt.Fscan(in, &root1)
			root1--
			fmt.Fprintln(out, hld.QuerySubtree(root1))
		}
	}

}
