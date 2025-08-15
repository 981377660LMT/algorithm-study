// CF176E Archaeology
// !https://www.luogu.com.cn/article/lhuhdi4k
// https://codeforces.com/problemset/problem/176/E
//
// 输入 n(1≤n≤1e5) 和一棵无向树的 n-1 条边（节点编号从 1 开始），每条边包含 3 个数 x y z(1≤z≤1e9)，表示有一条边权为 z 的边连接 x 和 y。
// 一开始有一个集合 S，初始为空。
// 然后输入 q(1≤q≤1e5) 和 q 个询问，格式如下：
// "+ v"：把点 v 加到集合 S 中，保证 v 不在 S 中。
// "- v"：把点 v 从集合 S 中删除，保证 v 在 S 中。
// "?"：输出包含 S 所有点的最小连通块（用最少的边连通 S 中所有点）的边权之和。
//
// 在 3553. 包含给定路径的最小带权子树 II 的基础上，我们需要动态增删点。
//
// 添加点时，我们需要知道 v 的 DFS 序最近的左右两个点 pre 和 nxt，这可以用平衡树 map 维护，key 是 DFN，value 是节点编号。
// 总和减去 dis(pre, nxt)，加上 dis(pre, v) 和 dis(v, nxt)。
// 删除则反过来。
// 输出答案时要除以 2，因为我们算的是回路。

package main

func main() {

}
