package main

import (
	"fmt"
)

type GridGraph struct {
	dims      []int
	dim       int
	volume    int
	partition []int
}

// 根据每个维度的大小创建一个网格图.
func NewGridGraph(dims ...int) *GridGraph {
	g := &GridGraph{
		dims: dims,
		dim:  len(dims),
	}
	r := []int{1}
	for i := len(dims) - 1; i >= 0; i-- {
		r = append(r, r[len(r)-1]*dims[i])
	}
	n := len(r)
	reversedR := make([]int, n)
	for i := 0; i < n; i++ {
		reversedR[i] = r[n-1-i]
	}
	g.volume = reversedR[0]
	g.partition = reversedR[1:]
	return g
}

// 将编号转换为多维坐标。
func (g *GridGraph) IdToPosition(id int) []int {
	if id < 0 || id >= g.volume {
		panic("编号超出范围")
	}
	pos := make([]int, g.dim)
	remain := id
	for i := 0; i < g.dim; i++ {
		pos[i] = remain / g.partition[i]
		remain = remain % g.partition[i]
	}
	return pos
}

// 将多维坐标转换为编号。
func (g *GridGraph) PositionToId(pos []int) int {
	if len(pos) != g.dim {
		panic("维度不匹配")
	}
	res := 0
	for i := 0; i < g.dim; i++ {
		if pos[i] < 0 || pos[i] >= g.dims[i] {
			panic(fmt.Sprintf("坐标 %d 超出范围", pos[i]))
		}
		res += g.partition[i] * pos[i]
	}
	return res
}

// 对编号 id 的所有邻居调用回调函数。
// 遍历所有维度，对每个维度，上下各一个邻居（如果存在）。
func (g *GridGraph) EnumerateNeightborIds(id int, f func(neighbor int)) {
	if id < 0 || id >= g.volume {
		panic("编号超出范围")
	}
	r := id
	for i := 0; i < g.dim; i++ {
		q := r / g.partition[i]
		r = r % g.partition[i]
		if q > 0 {
			f(id - g.partition[i])
		}
		if q < g.dims[i]-1 {
			f(id + g.partition[i])
		}
	}
}

// 对坐标 pos 的所有邻居调用回调函数。
// 回调函数共享同一个切片。
func (g *GridGraph) EnumerateNeighborPositions(pos []int, f func(neighbor []int)) {
	if len(pos) != g.dim {
		panic("维度不匹配")
	}
	pos = append(pos[:0:0], pos...)
	for i := 0; i < g.dim; i++ {
		if pos[i] > 0 {
			pos[i]--
			f(pos)
			pos[i]++
		}
		if pos[i] < g.dims[i]-1 {
			pos[i]++
			f(pos)
			pos[i]--
		}
	}
}

func main() {
	// 示例：创建一个 3x4 的二维网格
	grid := NewGridGraph(3, 4)
	fmt.Println("网格维度:", grid.dim)
	fmt.Println("网格总体积:", grid.volume)
	fmt.Println("每个维度大小:", grid.dims)
	fmt.Println("每个维度的 partition 权重:", grid.partition)

	// 将编号转换为位置
	N := 7
	pos := grid.IdToPosition(N)
	fmt.Printf("编号 %d 对应的位置: %v\n", N, pos)

	// 将位置转换为编号（验证与上面的编号是否一致）
	N2 := grid.PositionToId(pos)
	fmt.Printf("位置 %v 对应的编号: %d\n", pos, N2)

	// 输出编号 N 的邻居编号（在一维编号空间中）
	fmt.Printf("编号 %d 的邻居编号:\n", N)
	grid.EnumerateNeightborIds(N, func(neighbor int) {
		np := grid.IdToPosition(neighbor)
		fmt.Printf("%d -> 位置 %v\n", neighbor, np)
	})

	// 输出位置 pos 的邻居位置（直接操作多维坐标）
	fmt.Printf("位置 %v 的邻居位置:\n", pos)
	grid.EnumerateNeighborPositions(pos, func(neighbor []int) {
		num := grid.PositionToId(neighbor)
		fmt.Printf("%v -> 对应编号 %d\n", neighbor, num)
	})
}
