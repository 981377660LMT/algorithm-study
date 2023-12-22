// 稠密二分图最大匹配.
// O(n1^2*n2/64)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	Yuki421()
}

// https://yukicoder.me/problems/no/421
func Yuki421() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var H, W int
	fmt.Fscan(in, &H, &W)
	G := make([]string, H)
	for i := 0; i < H; i++ {
		fmt.Fscan(in, &G[i])
	}
	idx := make([][]int, H)
	for i := 0; i < H; i++ {
		idx[i] = make([]int, W)
	}
	a, b := 0, 0
	for x := 0; x < H; x++ {
		for y := 0; y < W; y++ {
			if (x+y)&1 == 0 {
				idx[x][y] = a
				a++
			}
			if (x+y)&1 == 1 {
				idx[x][y] = b
				b++
			}
		}
	}

	isIn := func(x, y int) bool {
		return 0 <= x && x < H && 0 <= y && y < W
	}

	dx := []int{1, 0, -1, 0, 1, 1, -1, -1}
	dy := []int{0, 1, 0, -1, 1, -1, 1, -1}

	adj := make([]*BS, a)
	for i := 0; i < a; i++ {
		adj[i] = NewBS(b, 0)
	}
	for x := 0; x < H; x++ {
		for y := 0; y < W; y++ {
			if (x+y)&1 == 1 {
				continue
			}
			for d := 0; d < 4; d++ {
				nx, ny := x+dx[d], y+dy[d]
				if !isIn(nx, ny) {
					continue
				}
				if G[x][y] == G[nx][ny] {
					continue
				}
				if G[x][y] == '.' {
					continue
				}
				if G[nx][ny] == '.' {
					continue
				}
				adj[idx[x][y]].Add(idx[nx][ny])
			}
		}
	}

	bm := NewBipartiteMatchingDense(a, b, adj)
	match := bm.MaxMatching()
	n := len(match)
	x, y := 0, 0
	for i := 0; i < H; i++ {
		for j := 0; j < W; j++ {
			if G[i][j] == 'w' {
				x++
			}
			if G[i][j] == 'b' {
				y++
			}
		}
	}
	x -= n
	y -= n
	res := 0
	res += 100 * n
	m := min(x, y)
	res += 10 * m
	res += x + y - 2*m
	fmt.Fprintln(out, res)
}

// 稠密二分图最大匹配.
type BipartiteMatchingDense struct {
	n1, n2         int32
	adjMatrix      []*BS
	match1, match2 []int32
	visited        *BS
}

func NewBipartiteMatchingDense(n1, n2 int, adjMatrix []*BS) *BipartiteMatchingDense {
	res := &BipartiteMatchingDense{n1: int32(n1), n2: int32(n2), adjMatrix: adjMatrix}
	res.match1 = make([]int32, n1)
	for i := 0; i < n1; i++ {
		res.match1[i] = -1
	}
	res.match2 = make([]int32, n2)
	for i := 0; i < n2; i++ {
		res.match2[i] = -1
	}
	res.visited = NewBS(n2, 1)
	for s := int32(0); s < int32(n1); s++ {
		res.bfs(s)
	}
	return res
}

func (bm *BipartiteMatchingDense) MaxMatching() (res [][2]int) {
	for v := int32(0); v < bm.n1; v++ {
		if bm.match1[v] != -1 {
			res = append(res, [2]int{int(v), int(bm.match1[v])})
		}
	}
	return
}

func (bm *BipartiteMatchingDense) MinVertexCover() (left, right []int) {
	queue := []int32{}
	bm.visited.Fill(1)
	done := make([]bool, bm.n1)
	for i := int32(0); i < bm.n1; i++ {
		if bm.match1[i] == -1 {
			done[i] = true
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		cand := bm.visited.And(bm.adjMatrix[cur])
		for v := cand.Next(0); v < int(bm.n2); v = cand.Next(v + 1) {
			bm.visited.Discard(v)
			to := bm.match2[v]
			if !done[to] {
				done[to] = true
				queue = append(queue, to)
			}
		}
	}
	for i := int32(0); i < bm.n1; i++ {
		if !done[i] {
			left = append(left, int(i))
		}
	}
	for i := int32(0); i < bm.n2; i++ {
		if !bm.visited.Has(int(i)) {
			right = append(right, int(i))
		}
	}
	return
}

func (bm *BipartiteMatchingDense) bfs(start int32) {
	if bm.match1[start] != -1 {
		return
	}
	queue := []int32{}
	prev := make([]int32, bm.n1)
	prev[start] = -1
	bm.visited.Fill(1)
	queue = append(queue, start)
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		cand := bm.visited.And(bm.adjMatrix[cur])
		for v := cand.Next(0); v < int(bm.n2); v = cand.Next(v + 1) {
			bm.visited.Discard(v)
			if bm.match2[v] != -1 {
				queue = append(queue, bm.match2[v])
				prev[bm.match2[v]] = cur
				continue
			}
			a, b := cur, int32(v)
			for a != -1 {
				t := bm.match1[a]
				bm.match1[a] = b
				bm.match2[b] = a
				a = prev[a]
				b = t
			}
			return
		}
	}
}

type BS struct {
	n    int
	data []uint64
}

func NewBS(n int, filledValue int) *BS {
	if !(filledValue == 0 || filledValue == 1) {
		panic("filledValue should be 0 or 1")
	}
	data := make([]uint64, (n+63)>>6)
	if filledValue == 1 {
		for i := range data {
			data[i] = ^uint64(0)
		}
		if n != 0 {
			data[len(data)-1] >>= (len(data) << 6) - n
		}
	}
	return &BS{n: n, data: data}
}

func (bs *BS) Add(i int) *BS {
	bs.data[i>>6] |= 1 << (i & 63)
	return bs
}

func (bs *BS) Has(i int) bool {
	return bs.data[i>>6]>>(i&63)&1 == 1
}

func (bs *BS) Discard(i int) {
	bs.data[i>>6] &^= 1 << (i & 63)
}

func (bs *BS) Fill(zeroOrOne int) {
	if zeroOrOne == 0 {
		for i := range bs.data {
			bs.data[i] = 0
		}
	} else {
		for i := range bs.data {
			bs.data[i] = ^uint64(0)
		}
		if bs.n != 0 {
			bs.data[len(bs.data)-1] >>= (len(bs.data) << 6) - bs.n
		}
	}
}

// 返回右侧第一个 1 的位置(`包含`当前位置).
//
//	如果不存在, 返回 n.
func (bs *BS) Next(index int) int {
	if index < 0 {
		index = 0
	}
	if index >= bs.n {
		return bs.n
	}
	k := index >> 6
	x := bs.data[k]
	s := index & 63
	x = (x >> s) << s
	if x != 0 {
		return (k << 6) | bs._lowbit(x)
	}
	for i := k + 1; i < len(bs.data); i++ {
		if bs.data[i] == 0 {
			continue
		}
		return (i << 6) | bs._lowbit(bs.data[i])
	}
	return bs.n
}

func (bs *BS) And(other *BS) *BS {
	res := NewBS(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] & v
	}
	return res
}

func (bs *BS) CopyAndResize(size int) *BS {
	newBits := make([]uint64, (size+63)>>6)
	copy(newBits, bs.data[:min(len(bs.data), len(newBits))])
	remainingBits := size & 63
	if remainingBits != 0 {
		mask := (1 << remainingBits) - 1
		newBits[len(newBits)-1] &= uint64(mask)
	}
	return &BS{data: newBits, n: size}
}

func (bs *BS) Resize(size int) {
	newBits := make([]uint64, (size+63)>>6)
	copy(newBits, bs.data[:min(len(bs.data), len(newBits))])
	remainingBits := size & 63
	if remainingBits != 0 {
		mask := (1 << remainingBits) - 1
		newBits[len(newBits)-1] &= uint64(mask)
	}
	bs.data = newBits
	bs.n = size
}

// 遍历所有 1 的位置.
func (bs *BS) ForEach(f func(pos int) (shouldBreak bool)) {
	for i, v := range bs.data {
		for ; v != 0; v &= v - 1 {
			j := (i << 6) | bs._lowbit(v)
			if f(j) {
				return
			}
		}
	}
}

func (bs *BS) String() string {
	sb := strings.Builder{}
	sb.WriteString("BS{")
	nums := []string{}
	bs.ForEach(func(pos int) bool {
		nums = append(nums, fmt.Sprintf("%d", pos))
		return false
	})
	sb.WriteString(strings.Join(nums, ","))
	sb.WriteString("}")
	return sb.String()
}

// (0, 1, 2, 3, 4) -> (-1, 0, 1, 1, 2)
func (bs *BS) _topbit(x uint64) int {
	if x == 0 {
		return -1
	}
	return 63 - bits.LeadingZeros64(x)
}

// (0, 1, 2, 3, 4) -> (-1, 0, 1, 0, 2)
func (bs *BS) _lowbit(x uint64) int {
	if x == 0 {
		return -1
	}
	return bits.TrailingZeros64(x)
}

func (bs *BS) _get(i int) uint64 {
	return bs.data[i>>6] >> (i & 63) & 1
}

func (bs *BS) _lastIndexOfOne() int {
	for i := len(bs.data) - 1; i >= 0; i-- {
		x := bs.data[i]
		if x != 0 {
			return (i << 6) | (bs._topbit(x))
		}
	}
	return -1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
