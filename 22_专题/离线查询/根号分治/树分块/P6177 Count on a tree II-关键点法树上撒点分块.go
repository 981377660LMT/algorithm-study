// https://www.luogu.com.cn/problem/P6177
// P6177 Count on a tree II-关键点法树上撒点分块
// 严格保证每个关键点到离它最近的祖先关键点的距离。
// !我们每次选择一个深度最大的非关键点，如果这个点的 1∼S级祖先都不是关键点，那么把它的 S级祖先设为关键点。由这个过程可知，距离不会超过 x。并且每标记一个关键点，至少有 S个点不会被标记。关键点数量也是对的。
// 给定一棵树，每个点带点权，每次询问给出 u,v，求 u到 v的简单路径上的不同权值个数。强制在线。
// https://mrsrz.github.io/lg6177/
// https://www.cnblogs.com/lyz09-blog/p/study-tree-block.html
package main

func main() {

}
