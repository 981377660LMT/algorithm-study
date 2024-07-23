package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"
)

type Random struct {
	seed     uint64
	hashBase uint64
}

func NewRandom() *Random { return &Random{seed: uint64(time.Now().UnixNano()/2 + 1)} }

func (r *Random) Rng() uint64 {
	r.seed ^= r.seed << 7
	r.seed ^= r.seed >> 9
	return r.seed
}
func (r *Random) Rng61() uint64 { return r.Rng() & ((1 << 61) - 1) }

// https://www.luogu.com.cn/problem/solution/CF1394B
// CF1394B Boboniu Walks on Graph
// 给定一个有向图，每个点的出度最大为k，k<=9.
// 每一条边都有一个权值，没有两条边有相同的权值
// !在图上走的时候，出度为 i 的边，只能走到边权大小为第 ci 小的边上面，
// !问c数组有多少种，要每个点都满足：从这个点开始走都可以回到这个点
//
// !只要按照枚举出来的 c 数组走，每个点的入度等于1就可以
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int32
	fmt.Fscan(in, &n, &m, &k)

	adjList := make([][][2]int32, n)
	for i := int32(0); i < m; i++ {
		var u, v, w int32
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		adjList[u] = append(adjList[u], [2]int32{v, w})
	}

	base := make([]uint64, n)
	random := NewRandom()
	full := uint64(0)
	for i := int32(0); i < n; i++ {
		base[i] = random.Rng61()
		full += base[i]
	}

	// dat[k][i] 表示第 k 个点的第 i 小的边的权值
	data := make([][]uint64, k+1)
	for i := int32(0); i <= k; i++ {
		data[i] = make([]uint64, k)
	}
	for v := int32(0); v < n; v++ {
		tmp := make([][2]int32, 0, len(adjList[v]))
		for _, to := range adjList[v] {
			tmp = append(tmp, [2]int32{to[1], to[0]}) // (cost, to)
		}
		sort.Slice(tmp, func(i, j int) bool { return tmp[i][0] < tmp[j][0] })
		for i, e := range tmp {
			data[len(tmp)][i] += base[e[1]]
		}
	}

	res := 0
	var dfs func(int32, uint64)
	dfs = func(index int32, now uint64) {
		if index == k+1 {
			if now == full {
				res++
			}
			return
		}
		for i := int32(0); i < index; i++ {
			dfs(index+1, now+data[index][i])
		}
	}
	dfs(1, 0)

	fmt.Fprintln(out, res)
}
