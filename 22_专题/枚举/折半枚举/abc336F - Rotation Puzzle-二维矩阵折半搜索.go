// https://atcoder.jp/contests/abc336/tasks/abc336_f
// 折半搜索/双向bfs
// 给定一个n*m(3<=n,m<=8)的网格,0到n*m-1每个数字恰好出现一次.
// !现在可以以(0,0),(1,0),(0,1)或(1,1)为左上角180度顺时针旋转(n-1)*(m-1)大小的矩阵.
// 最多操作`20`次,问能否使得网格变成有序的,即0-n*m-1从左到右从上到下依次排列.
// 如果能,输出最少操作次数,否则输出-1.

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var H, W int
	fmt.Fscan(in, &H, &W)
	grid := make([][]int, H)
	for i := 0; i < H; i++ {
		grid[i] = make([]int, W)
		for j := 0; j < W; j++ {
			fmt.Fscan(in, &grid[i][j])
			grid[i][j]--
		}
	}

	fmt.Fprintln(out, RotationPuzzle(grid))
}

const INF int = 1e18

func RotationPuzzle(grid [][]int) int {
	CENTER_4 := [][]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}

	ROW, COL := len(grid), len(grid[0])
	target := make([][]int, ROW)
	for i := 0; i < ROW; i++ {
		target[i] = make([]int, COL)
		for j := 0; j < COL; j++ {
			target[i][j] = i*COL + j
		}
	}
	hashBase := NewRandom().GetHashBase2D(grid)

	// 以(leftX, leftY)为左上角顺时针旋转180度.
	rotate := func(g [][]int, leftX, leftY int) [][]int {
		res := make([][]int, ROW)
		for i := 0; i < ROW; i++ {
			res[i] = make([]int, COL)
			for j := 0; j < COL; j++ {
				res[i][j] = g[i][j]
			}
		}
		for i := 0; i < ROW-1; i++ {
			for j := 0; j < COL-1; j++ {
				res[leftX+i][leftY+j] = g[leftX+ROW-2-i][leftY+COL-2-j]
			}
		}
		return res
	}

	getStates := func(g [][]int) map[uint64]int {
		res := make(map[uint64]int)
		var dfs func(curGrid [][]int, step int)
		dfs = func(curGrid [][]int, step int) {
			matHash := uint64(0)
			for i := 0; i < ROW; i++ {
				for j := 0; j < COL; j++ {
					matHash += uint64(curGrid[i][j]) * hashBase[i][j]
				}
			}
			if _, ok := res[matHash]; ok {
				if res[matHash] <= step {
					return
				}
			}
			res[matHash] = step
			if step == 10 {
				return
			}
			for _, v := range CENTER_4 {
				dfs(rotate(curGrid, v[0], v[1]), step+1)
			}
		}
		dfs(g, 0)
		return res
	}

	mp1, mp2 := getStates(grid), getStates(target)
	res := INF
	for k, v := range mp1 {
		if _, ok := mp2[k]; ok {
			res = min(res, v+mp2[k])
		}
	}
	if res == INF {
		return -1
	}
	return res
}

type Random struct {
	seed     uint64
	hashBase uint64
}

func NewRandom() *Random                 { return &Random{seed: uint64(time.Now().UnixNano()/2 + 1)} }
func NewRandomWithSeed(seed int) *Random { return &Random{seed: uint64(seed)} }

func (r *Random) Rng() uint64 {
	r.seed ^= r.seed << 7
	r.seed ^= r.seed >> 9
	return r.seed
}

func (r *Random) Next() uint64 { return r.Rng() }

func (r *Random) RngWithMod(mod int) uint64 { return r.Rng() % uint64(mod) }

// [left, right]
func (r *Random) RandInt(min, max int) uint64 { return uint64(min) + r.Rng()%(uint64(max-min+1)) }

// [start:stop:step]
func (r *Random) RandRange(start, stop int, step int) uint64 {
	width := stop - start
	// Fast path.
	if step == 1 {
		return uint64(start) + r.Rng()%uint64(width)
	}
	var n uint64
	if step > 0 {
		n = uint64((width + step - 1) / step)
	} else {
		n = uint64((width + step + 1) / step)
	}
	return uint64(start) + uint64(step)*(r.Rng()%n)
}

// FastShuffle
func (r *Random) Shuffle(nums []int) {
	for i := range nums {
		rand := r.RandInt(0, i)
		nums[i], nums[rand] = nums[rand], nums[i]
	}
}

func (r *Random) Sample(nums []int, k int) []int {
	nums = append(nums[:0:0], nums...)
	r.Shuffle(nums)
	return nums[:k]
}

// 元组哈希
func (r *Random) HashPair(a, b int) uint64 {
	if r.hashBase == 0 {
		r.hashBase = r.Rng()
	}
	return uint64(a)*r.hashBase + uint64(b)
}

func (r *Random) GetHashBase1D(nums []int) []uint64 {
	hashBase := make([]uint64, len(nums))
	for i := range hashBase {
		hashBase[i] = r.Rng()
	}
	return hashBase
}

func (r *Random) GetHashBase2D(nums [][]int) [][]uint64 {
	hashBase := make([][]uint64, len(nums))
	for i := range hashBase {
		hashBase[i] = make([]uint64, len(nums[i]))
		for j := range hashBase[i] {
			hashBase[i][j] = r.Rng()
		}
	}
	return hashBase
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
