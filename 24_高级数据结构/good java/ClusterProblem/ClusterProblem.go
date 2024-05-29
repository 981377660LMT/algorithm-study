// ClusterProblem 是一个聚类问题的解决方案。聚类是一种无监督学习的方法，用于将数据集中的对象分成不同的组或簇，使得同一组内的对象相似度较高，而不同组之间的相似度较低。
// 该代码实现了一种基于距离矩阵的聚类算法。给定一个距离矩阵，其中 dists[i * n + j] 表示第 i 个对象和第 j 个对象之间的距离（距离大于0），并且 dists[i * n + i] = 0。
// 算法的目标是将对象分成不同的簇，使得同一簇内的对象之间的距离较小。
// 该算法的主要步骤如下：
// 1. 根据距离矩阵的大小创建一个索引数组 indices，用于存储对象之间的索引关系。
// 2. 使用快速排序对索引数组 indices 进行排序，排序的依据是根据距离矩阵中的距离进行比较。
// 3. 初始化一个 DSUExt 对象，DSUExt 是一个扩展了并查集（Disjoint Set Union）的数据结构，用于管理聚类过程中的连接关系。
// 4. 遍历排序后的索引数组，根据对象之间的距离进行合并操作，将相似的对象连接到同一个簇中。
// 5. 根据合并后的连接关系，构建每个簇的布尔数组，表示该簇中的对象。
// 6. 返回每个簇的布尔数组，即为聚类的结果。
// 该算法的时间复杂度为 O(n^2)，其中 n 是对象的数量。它通过合并相似的对象来构建聚类结果，适用于基于距离矩阵的聚类问题。

package main

import (
	"fmt"
	"sort"
)

func main() {
	// 示例用法
	n := int32(3)
	dist := func(i, j int32) int {
		return int(i + j + 1)
	}
	res := ClusterProblem(n, dist)
	fmt.Println(res) // Output: [false true true true]
}

func ClusterProblem(n int32, dist func(i, j int32) int) []bool {
	dists := make([]int, n*n)
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < n; j++ {
			if i != j {
				dists[i*n+j] = dist(i, j)
			}
		}
	}
	indices := make([]int32, n*n)
	for i := int32(0); i < n*n; i++ {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return dists[indices[i]] < dists[indices[j]]
	})
	dsu := NewDSUExt(n)
	m := int32(len(indices))
	for i := int32(0); i < m; i++ {
		j := i
		for j+1 < m && dists[indices[i]] == dists[indices[j+1]] {
			j++
		}
		for k := i; k <= j; k++ {
			a := indices[k] / n
			b := indices[k] - a*n
			dsu.Union(a, b)
		}
		for k := i; k <= j; k++ {
			a := indices[k] / n
			dsu.AddEdge(dsu.Find(a))
		}
		i = j
	}
	return dsu.dp[dsu.Find(0)]
}

type DSUExt struct {
	fa   []int32
	size []int32
	dp   [][]bool
	edge []int32
}

func NewDSUExt(n int32) *DSUExt {
	fa := make([]int32, n)
	size := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fa[i] = i
		size[i] = 1
	}
	dp := make([][]bool, n)
	for i := int32(0); i < n; i++ {
		dp[i] = make([]bool, n+1)
	}
	edge := make([]int32, n)
	return &DSUExt{fa: fa, size: size, dp: dp, edge: edge}
}

func (d *DSUExt) AddEdge(x int32) {
	d.edge[x]++
	if d.edge[x] == d.size[x]*d.size[x] {
		d.dp[x][1] = true
	}
}

func (d *DSUExt) Union(a, b int32) {
	a, b = d.Find(a), d.Find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}

	// preMerge
	{
		for i := d.size[a]; i >= 1; i-- {
			if !d.dp[a][i] {
				continue
			}
			d.dp[a][i] = false
			for j := int32(1); j <= d.size[b]; j++ {
				if d.dp[b][j] {
					d.dp[a][i+j] = true
				}
			}
		}
		d.edge[a] += d.edge[b]
	}

	d.size[a] += d.size[b]
	d.fa[b] = a
}

func (d *DSUExt) Find(i int32) int32 {
	if d.fa[d.fa[i]] == d.fa[i] {
		return d.fa[i]
	}
	d.fa[i] = d.Find(d.fa[i])
	return d.fa[i]
}
