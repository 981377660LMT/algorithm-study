// https://www.luogu.com.cn/problem/CF938G
// 每条无向边有边权和一段作用时间区间，查询两个点之间的异或最短路。
//
// 首先考虑怎么求两个点之间的异或最短路，其实就是P4151 最大XOR和路径
// 考虑如果一颗树的话，之间的答案就是路径上所有边异或起来的值。
// 那么推广到图上去，不难想到就是多了一些环，而且这个环走一圈是没贡献的（异或为零）。所以如果你要从环的一段走向另一端，实际上只有两种走法。
//
// 题意就转化为选异或任意环的权值(cycle基底)，使异或值最大。
//
// !代码实现见 `线段树分治` 下的习题.

package main
