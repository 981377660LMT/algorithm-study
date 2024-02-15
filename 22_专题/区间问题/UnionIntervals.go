// IntervalUnion/UnionInterval/IntervalGraphUnion
// 区间合并/合并区间

package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	fmt.Println(UnionIntervals([][2]int{{1, 3}, {3, 4}, {5, 7}, {6, 8}}))
}

// 给定n个区间[start,end), 返回合并后的并查集.
func UnionIntervals(intervals [][2]int) *UnionFindArray {
	n := len(intervals)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		a, b := intervals[order[i]], intervals[order[j]]
		if a[0] == b[0] {
			return a[1] > b[1]
		}
		return a[0] < b[0]
	})

	uf := NewUnionFindArray(n)
	keep := make([]int, 0, n)
	for _, j := range order {
		if len(keep) > 0 {
			i := keep[len(keep)-1]
			startI, endI := intervals[i][0], intervals[i][1]
			startJ, endJ := intervals[j][0], intervals[j][1]
			if endJ <= endI && endJ-startJ < endI-startI {
				uf.Union(i, j)
				continue
			}
		}
		keep = append(keep, j)
	}

	for k := 0; k < len(keep)-1; k++ {
		i, j := keep[k], keep[k+1]
		startI, endI := intervals[i][0], intervals[i][1]
		startJ, endJ := intervals[j][0], intervals[j][1]
		if max(startI, startJ) < min(endI, endJ) {
			uf.Union(i, j)
		}
	}

	return uf
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int
	n    int
	data []int
}

func NewUnionFindArray(n int) *UnionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{
		Part: n,
		n:    n,
		data: data,
	}
}

// 按秩合并.
func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	ufa.Part--
	if cb != nil {
		cb(root1, root2)
	}
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetSize(key int) int {
	return -ufa.data[ufa.Find(key)]
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	groups := ufa.GetGroups()
	keys := make([]int, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
