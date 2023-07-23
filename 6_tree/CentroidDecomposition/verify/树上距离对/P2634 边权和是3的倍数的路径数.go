// https://www.luogu.com.cn/problem/P2634

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://www.luogu.com.cn/problem/P2634
	// !有多少条路径边权和是3的倍数
	// !n<=2e4 O(n*logn*logn)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	g := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		g[u] = append(g[u], Edge{to: v, weight: w})
		g[v] = append(g[v], Edge{to: u, weight: w})
	}

	// 全局状态
	centTree, root := CentroidDecomposition(g)
	removed := make([]bool, n)
	res := 0

	var collectDist func(cur, pre, distMod int, mp map[int]int)
	collectDist = func(cur, pre, distMod int, mp map[int]int) {
		mp[distMod]++
		for _, e := range g[cur] {
			next, cost := e.to, e.weight
			if next != pre && !removed[next] {
				collectDist(next, cur, (distMod+cost)%3, mp)
			}
		}
	}

	var decomposition func(cur, pre int)
	decomposition = func(cur, pre int) {
		removed[cur] = true
		for _, next := range centTree[cur] { // 点分树的子树中的答案(不经过重心)
			if !removed[next] {
				decomposition(next, cur)
			}
		}
		removed[cur] = false

		// 统计经过重心的路径
		counter := map[int]int{0: 1} // mod3 -> count
		for _, e := range g[cur] {
			next, cost := e.to, e.weight
			if next == pre || removed[next] {
				continue
			}

			sub := map[int]int{}
			collectDist(next, cur, cost%3, sub)
			for k, v := range sub {
				res += v * counter[(3-k)%3]
			}
			for k, v := range sub {
				counter[k] += v // !注意这里不能和上面的循环合并
			}
		}
	}

	decomposition(root, -1)
	res = res*2 + n
	gcd := func(a, b int) int {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}
	gcd_ := gcd(res, n*n)
	fmt.Fprintf(out, "%d/%d", res/gcd_, n*n/gcd_)
}

type Edge = struct{ to, weight int }

// 树的重心分解, 返回点分树和点分树的根
//
//	!tree: `无向`树的邻接表.
//	centTree: 重心互相连接形成的有根树, 可以想象把树拎起来, 重心在树的中心，连接着各个子树的重心...
//	root: 点分树的根
func CentroidDecomposition(tree [][]Edge) (centTree [][]int, root int) {
	n := len(tree)
	subSize := make([]int, n)
	removed := make([]bool, n)
	centTree = make([][]int, n)

	var getSize func(cur, parent int) int
	var getCentroid func(cur, parent, mid int) int
	var build func(cur int) int
	getSize = func(cur, parent int) int {
		subSize[cur] = 1
		for _, e := range tree[cur] {
			next := e.to
			if next == parent || removed[next] {
				continue
			}
			subSize[cur] += getSize(next, cur)
		}
		return subSize[cur]
	}
	getCentroid = func(cur, parent, mid int) int {
		for _, e := range tree[cur] {
			next := e.to
			if next == parent || removed[next] {
				continue
			}
			if subSize[next] > mid {
				return getCentroid(next, cur, mid)
			}
		}
		return cur
	}
	build = func(cur int) int {
		centroid := getCentroid(cur, -1, getSize(cur, -1)/2)
		removed[centroid] = true
		for _, e := range tree[centroid] {
			next := e.to
			if !removed[next] {
				centTree[centroid] = append(centTree[centroid], build(next))
			}
		}
		removed[centroid] = false
		return centroid
	}

	root = build(0)
	return
}
