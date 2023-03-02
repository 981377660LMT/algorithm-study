// https://www.luogu.com.cn/problem/CF321C
// Fox Ciel需要在每个城市分配一个官员。每个官员都有一个等级---一个’A’到’Z’之间的字母。
// !所以会有26个不同的等级，’A’是最高的，’Z’是最低的。
// 如果x和y是两个不同的城市并且他们的官员拥有`相同的等级`，
// !那么x和y之间的简单路径中一定存在一个城市z有更高等级的官员
// 这个规则可以保证两个`同等级官员之间的通信`会由较高等级的官员监控。
// 帮助Ciel制定一个有效的计划，如果这是不可能的，输出"Impossible!"。
// n<=1e5
// 我们要节省使用的等级数量，就要保证分解问题的次数尽可能少，也就是说，要让每次切出来的子树都尽可能小
// 这是重心的性质。这保证了我们分解问题的次数不会超过logn
// !在点分树中,根节点最大,子树递减即可.

// 点分治模型

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
	g := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		g[u] = append(g[u], Edge{to: v, cost: 1})
		g[v] = append(g[v], Edge{to: u, cost: 1})
	}

	// 全局状态
	centTree, root := CentroidDecomposition(g)
	res := make([]int, n)
	var dfs func(cur, assign int)
	dfs = func(cur, assign int) {
		res[cur] = assign
		for _, next := range centTree[cur] {
			dfs(next, assign+1)
		}
	}
	dfs(root, 0)

	for i := 0; i < n; i++ {
		fmt.Print(string(rune('A'+res[i])), " ")
	}
}

type Edge = struct{ to, cost int }

// 树的重心分解, 返回点分树和点分树的根
//  g: 原图
//  centTree: 重心互相连接形成的有根树, 可以想象把树拎起来, 重心在树的中心，连接着各个子树的重心...
//  root: 点分树的根
func CentroidDecomposition(g [][]Edge) (centTree [][]int, root int) {
	n := len(g)
	subSize := make([]int, n)
	removed := make([]bool, n)
	centTree = make([][]int, n)

	var getSize func(cur, parent int) int
	var getCentroid func(cur, parent, mid int) int
	var build func(cur int) int
	getSize = func(cur, parent int) int {
		subSize[cur] = 1
		for _, e := range g[cur] {
			next := e.to
			if next == parent || removed[next] {
				continue
			}
			subSize[cur] += getSize(next, cur)
		}
		return subSize[cur]
	}
	getCentroid = func(cur, parent, mid int) int {
		for _, e := range g[cur] {
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
		for _, e := range g[centroid] {
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
