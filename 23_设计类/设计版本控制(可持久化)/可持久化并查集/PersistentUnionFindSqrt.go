// Persistence Union Find
//
// Description:
//   Use persistent array instead of standard array in union find data structure
//
// Complexity:
//   O(a* int(n)), where int(n) is a complexity of persistent array

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

func main() {
	yosupo()
}

func demo() {
	n := int(1e5)
	uf := NewPersistentUnionFindSqrt(n)
	time1 := time.Now()
	for i := 0; i < n-1; i++ {
		uf = uf.Union(i, i+1)
	}
	fmt.Println(time.Since(time1))
}

func yosupo() {
	// https://judge.yosupo.jp/submission/130167
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	git := make([]*PersistentUnionFindSqrt, 0, q+1)
	uf := NewPersistentUnionFindSqrt(n)
	git = append(git, uf)

	for i := 0; i < q; i++ {
		var op, version, u, v int
		fmt.Fscan(in, &op, &version, &u, &v)
		version++
		root := git[version]
		if op == 0 {
			newRoot := root.Union(u, v)
			root = newRoot
		} else {
			if root.IsConnected(u, v) {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 0)
			}
		}
		git = append(git, root)
	}
}

// 可持久化并查集.
type PersistentUnionFindSqrt struct {
	parent *_PS
}

func NewPersistentUnionFindSqrt(n int) *PersistentUnionFindSqrt {
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = -1
	}
	return &PersistentUnionFindSqrt{parent: _NPS(parent)}
}

func (uf *PersistentUnionFindSqrt) Union(u, v int) *PersistentUnionFindSqrt {
	root1 := uf.Find(u)
	root2 := uf.Find(v)
	if root1 == root2 {
		return uf
	}
	p1 := uf.parent.Get(root1)
	p2 := uf.parent.Get(root2)
	if p1 > p2 {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	newP := uf.parent.Set(root1, p1+p2)
	newP = newP.Set(root2, root1)
	return &PersistentUnionFindSqrt{parent: newP}
}

func (uf *PersistentUnionFindSqrt) Find(u int) int {
	for {
		p := uf.parent.Get(u)
		if p < 0 {
			break
		}
		u = p
	}
	return u
}

func (uf *PersistentUnionFindSqrt) IsConnected(u, v int) bool {
	return uf.Find(u) == uf.Find(v)
}

func (uf *PersistentUnionFindSqrt) GetSize(u int) int {
	return -uf.parent.Get(uf.Find(u))
}

func GetGroups(n int, uf *PersistentUnionFindSqrt) map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < n; i++ {
		groups[uf.Find(i)] = append(groups[uf.Find(i)], i)
	}
	return groups
}

type _PS struct {
	arr     []int
	opIndex []int
	opValue []int
	opLen   int
}

func _NPS(arr []int) *_PS {
	sqrt := int(math.Sqrt(float64(len(arr)))) + 1
	return &_PS{arr: arr, opIndex: make([]int, 0, sqrt), opValue: make([]int, 0, sqrt)}
}

func (sa *_PS) Get(i int) int {
	for j := sa.opLen - 1; j >= 0; j-- {
		if sa.opIndex[j] == i {
			return sa.opValue[j]
		}
	}
	return sa.arr[i]
}

func (sa *_PS) Set(i, v int) *_PS {
	sa.opIndex = append(sa.opIndex, i)
	sa.opValue = append(sa.opValue, v)
	n := len(sa.arr)
	if tmp := len(sa.opIndex); tmp*tmp <= n {
		return &_PS{arr: sa.arr, opIndex: sa.opIndex, opValue: sa.opValue, opLen: tmp}
	}
	newArr := make([]int, n)
	copy(newArr, sa.arr)
	for i := range sa.opIndex {
		newArr[sa.opIndex[i]] = sa.opValue[i]
	}
	return _NPS(newArr)
}
