// 没有js的版本快

package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

func main() {
	arr := NewSqrtArray(make([]int, 20))
	rands := make([]int, 4e5)
	for i := 0; i < 4e5; i++ {
		rands[i] = rand.Intn(i + 10)
	}
	time1 := time.Now()
	for i := 0; i < 4e5; i++ {
		arr.Insert(rands[i], i)
	}
	fmt.Println(time.Since(time1))
}

type T = int
type SqrtArray struct {
	n int
	x [][]T
}

func NewSqrtArray(nums []T) *SqrtArray {
	res := &SqrtArray{}
	n := len(nums)
	if n == 0 {
		return res
	}

	bCount := int(math.Sqrt(float64(n)))
	bSize := (n + bCount - 1) / bCount
	newB := make([][]T, bCount)
	for i := 0; i < bCount; i++ {
		newB[i] = []T{}
		for j := i * bSize; j < min((i+1)*bSize, n); j++ {
			newB[i] = append(newB[i], nums[j])
		}
	}
	res.n = n
	res.x = newB
	return res
}

func (sa *SqrtArray) Set(i int, v T) {
	if i == sa.n-1 {
		bi := len(sa.x) - 1
		sa.x[bi][len(sa.x[bi])-1] = v
		return
	}
	bi := 0
	for ; i >= len(sa.x[bi]); i, bi = i-len(sa.x[bi]), bi+1 {
	}
	sa.x[bi][i] = v
}

func (sa *SqrtArray) Get(i int) T {
	if i == sa.n-1 {
		bi := len(sa.x) - 1
		return sa.x[bi][len(sa.x[bi])-1]
	}
	j := 0
	for ; i >= len(sa.x[j]); i, j = i-len(sa.x[j]), j+1 {
	}
	return sa.x[j][i]
}

func (sa *SqrtArray) Append(v T) {
	sa.Insert(sa.n, v)
}

func (sa *SqrtArray) Pop() T {
	return sa.Erase(sa.n - 1)
}

func (sa *SqrtArray) AppendLeft(v T) {
	sa.Insert(0, v)
}

func (sa *SqrtArray) PopLeft() T {
	return sa.Erase(0)
}

func (sa *SqrtArray) Erase(i int) (res T) {
	bi := 0
	if i == sa.n-1 {
		bi = len(sa.x) - 1
		res = sa.x[bi][len(sa.x[bi])-1]
		sa.x[bi] = sa.x[bi][:len(sa.x[bi])-1]
	} else {
		for ; bi >= len(sa.x[bi]); i, bi = i-len(sa.x[bi]), bi+1 {
		}
		res = sa.x[bi][i]
		sa.x[bi] = append(sa.x[bi][:i], sa.x[bi][i+1:]...)
	}
	sa.n--

	if len(sa.x[bi]) == 0 {
		sa.x = append(sa.x[:bi], sa.x[bi+1:]...)
	}
	return res
}

func (sa *SqrtArray) Insert(i int, v T) {
	if sa.n == 0 {
		sa.x = append(sa.x, []T{})
	}

	bi := 0
	if i >= sa.n {
		bi = len(sa.x) - 1
		sa.x[bi] = append(sa.x[bi], v)
	} else {
		for ; bi < len(sa.x) && i >= len(sa.x[bi]); i, bi = i-len(sa.x[bi]), bi+1 {
		}
		sa.x[bi] = append(sa.x[bi][:i], append([]T{v}, sa.x[bi][i:]...)...)
	}
	sa.n++

	sqrt := int(math.Sqrt(float64(sa.n)))
	if len(sa.x[bi]) > 2*sqrt {
		y := sa.x[bi][sqrt:]
		sa.x[bi] = sa.x[bi][:sqrt]
		sa.x = append(sa.x[:bi+1], append([][]T{y}, sa.x[bi+1:]...)...)
	}
}

func (sa *SqrtArray) Len() int {
	return sa.n

}

func (sa *SqrtArray) ForEach(f func(value T, index int)) {
	ptr := 0
	for i := 0; i < len(sa.x); i++ {
		for j := 0; j < len(sa.x[i]); j++ {
			f(sa.x[i][j], ptr)
			ptr++
		}
	}
}

func (sa *SqrtArray) String() string {
	sb := make([]string, 0, sa.n)
	sa.ForEach(func(value T, _ int) {
		sb = append(sb, fmt.Sprintf("%d", value))
	})
	return fmt.Sprintf("SqrtArray[%s]", strings.Join(sb, ", "))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
