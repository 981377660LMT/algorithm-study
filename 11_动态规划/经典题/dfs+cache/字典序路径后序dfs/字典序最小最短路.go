// 字典序最小最短路
// 入门经典第二版 p.173
// 理想路径（NEERC10）https://codeforces.com/gym/101309 I 题
// 边的字典序最小:
//   从终点倒着 BFS 求最短路，然后从起点开始一层一层向终点走，每一步都选颜色最小的，并记录最小颜色对应的所有节点，供下一层遍历
// 点的字典序最小:
//   每一步需选择符合 dis[w] == dis[v]-1 的下标最小的顶点
// LC499 https://leetcode.cn/problems/the-maze-iii/
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/graph.go
package main

const INF int = 1e18

func lexicographicallySmallestShortestPath(g [][]struct{ to, color int }, st, end int) []int {
	dis := make([]int, len(g))
	for i := range dis {
		dis[i] = INF
	}
	dis[end] = 0
	q := []int{end}
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, e := range g[v] {
			if w := e.to; dis[v]+1 < dis[w] {
				dis[w] = dis[v] + 1
				q = append(q, w)
			}
		}
	}

	if dis[st] == INF {
		return nil
	}

	res := []int{}
	check := []int{st}
	inC := make([]bool, len(g))
	inC[st] = true
	for loop := dis[st]; loop > 0; loop-- {
		minC := INF
		tmp := check
		check = nil
		for _, v := range tmp {
			for _, e := range g[v] {
				if w, c := e.to, e.color; dis[w] == dis[v]-1 {
					if c < minC {
						for _, w := range check {
							inC[w] = false
						}
						minC, check, inC[w] = c, []int{w}, true
					} else if c == minC && !inC[w] {
						check = append(check, w)
						inC[w] = true
					}
				}
			}
		}
		res = append(res, minC)
	}
	return res
}
