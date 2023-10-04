// P3603 雪辉-树分块+bitset
// https://www.luogu.com.cn/problem/P3603
// https://oi-wiki.org/ds/tree-decompose/
// 给你一棵 n 个节点且带点权的树，m 个询问，
// 每个询问给你多条链，请你输出这几条链的点的集合并的颜色数和 mex。

// 方法1：树上撒点分块法(预处理出互为祖孙关系的关键点之间的数的信息)
// 方法2：跳重链分块(路径上的若干条重链的 bitset 并起来)
// 重链上的点的 dfn 序是连续的，序列分块即可
package main

func main() {

}
