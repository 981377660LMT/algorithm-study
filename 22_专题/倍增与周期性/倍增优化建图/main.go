// 倍增优化建图
// https://taodaling.github.io/blog/2020/03/18/binary-lifting/

package main

func main() {

}

// P5344 【XR-1】逛森林
// https://www.luogu.com.cn/problem/P5344
// 也可以隐式建图.

// P9520 [JOISC2022] 监狱
// https://www.luogu.com.cn/problem/P9520
// 对于 n个点的树，有 m条起点与终点各不相同的行进路线形如 si→ti，
// 允许从某个点移动至相邻点，问能否在不存在某个点所在人数 >1的情况下完成所有行进路线。

// Beautiful Tree
// https://www.luogu.com.cn/problem/CF1904F
// 给出一棵树，与 m 条限制，每条限制为一条路径上点权最大/小的点的编号固定。
// 请你为图分配 1∼n 的点权使得满足所有限制。
// 限制可以看成规定点点权大/于路径上的其它点，我们把 a 的点权大于 b 的点权的限制视作一个有向边，则有解当且仅当没有环，拓扑排序分配即可。
// !树剖 + 线段树优化建图O(nlog^2)，可以倍增优化成 O(nlogn)。
