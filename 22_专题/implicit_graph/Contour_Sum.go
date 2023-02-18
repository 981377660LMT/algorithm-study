// #include "ds/fenwicktree/fenwicktree.hpp"

// // 点加算、距離区間での和
// template <typename GT, typename AbelGroup>
// struct Contour_Sum {
//   int N;
//   GT& G;
//   using X = typename AbelGroup::value_type;
//   FenwickTree<AbelGroup> bit;
//   // centroid ごと、方向ごと
//   vvc<int> bit_range;
//   // 方向ラベル、重心からの距離、bit でのindex
//   vvc<tuple<int, int, int>> dat;

//   Contour_Sum(GT& G) : N(G.N), G(G) {
//     assert(!G.is_directed());
//     vc<X> v_vals(N, AbelGroup::unit());
//     build(v_vals);
//   }
//   Contour_Sum(GT& G, const vc<X>& v_vals) : N(G.N), G(G) {
//     assert(!G.is_directed());
//     build(v_vals);
//   }

//   void add(int v, X val) {
//     for (auto&& [k, x, i]: dat[v]) bit.add(i, val);
//   }

//   // v を中心として、距離 [l, r) の範囲の和
//   X sum(int v, int l, int r) {
//     X sm = AbelGroup::unit();
//     for (auto [k, x, i]: dat[v]) {
//       int lo = l - x, hi = r - x;
//       int K = k;
//       if (k < 0) { lo -= 2, hi -= 2, K = ~k; }
//       int n = len(bit_range[K]) - 2;
//       chmax(lo, 0);
//       chmin(hi, n + 1);
//       if (lo >= hi) continue;
//       int a = bit_range[K][lo], b = bit_range[K][hi];
//       X val = bit.prod(a, b);
//       if (k < 0) { val = AbelGroup::inverse(val); }
//       sm = AbelGroup::op(sm, val);
//     }
//     return sm;
//   }

//   void build(const vc<X>& v_vals) {
//     int nxt_bit_idx = 0;
//     vc<int> done(N, 0);
//     vc<int> sz(N, 0);
//     vc<int> par(N, -1);
//     vc<int> dist(N, -1);
//     vc<pair<int, int>> st;
//     bit_range.resize(N);
//     dat.resize(N);
//     st.eb(0, N);

//     while (len(st)) {
//       int v0 = st.back().fi;
//       int n = st.back().se;
//       st.pop_back();
//       int c = -1;
//       {
//         auto dfs = [&](auto& dfs, int v) -> int {
//           sz[v] = 1;
//           for (auto&& e: G[v])
//             if (e.to != par[v] && !done[e.to]) {
//               par[e.to] = v;
//               sz[v] += dfs(dfs, e.to);
//             }
//           if (c == -1 && n - sz[v] <= n / 2) { c = v; }
//           return sz[v];
//         };
//         dfs(dfs, v0);
//       }
//       // center からの bfs。部分木サイズもとっておく。
//       done[c] = 1;
//       {
//         int off = nxt_bit_idx;
//         vc<int> que;
//         auto add = [&](int v, int d, int p) -> void {
//           if (dist[v] != -1) return;
//           sz[v] = 1;
//           dist[v] = d;
//           par[v] = p;
//           que.eb(v);
//         };
//         int p = 0;
//         add(c, 0, -1);
//         while (p < len(que)) {
//           auto v = que[p++];
//           for (auto&& e: G[v]) {
//             if (done[e.to]) continue;
//             add(e.to, dist[v] + 1, v);
//           }
//         }
//         FOR_R(i, 1, len(que)) {
//           int v = que[i];
//           sz[par[v]] += sz[v];
//         }
//         // 距離ごとのカウント
//         int max_d = dist[que.back()];
//         vc<int> count(max_d + 1);
//         // 重心、方向ラベル、重心からの距離、bit でのindex
//         for (auto&& v: que) {
//           dat[v].eb(c, dist[v], nxt_bit_idx++);
//           count[dist[v]]++;
//           par[v] = -1;
//           dist[v] = -1;
//         }
//         bit_range[c] = cumsum<int, int>(count);
//         for (auto&& x: bit_range[c]) x += off;
//       }
//       // 方向ごとの bfs
//       for (auto&& e: G[c]) {
//         int off = nxt_bit_idx;
//         int nbd = e.to;
//         if (done[nbd]) continue;
//         int K = len(bit_range);
//         vc<int> que;
//         auto add = [&](int v, int d) -> void {
//           if (dist[v] != -1 || v == c) return;
//           dist[v] = d;
//           que.eb(v);
//         };
//         int p = 0;
//         add(nbd, 0);
//         while (p < len(que)) {
//           auto v = que[p++];
//           for (auto&& e: G[v]) {
//             if (done[e.to]) continue;
//             add(e.to, dist[v] + 1);
//           }
//         }
//         // 距離ごとのカウント
//         int max_d = dist[que.back()];
//         vc<int> count(max_d + 1);
//         // 重心、方向ラベル、重心からの距離、bit でのindex
//         for (auto&& v: que) {
//           dat[v].eb(~K, dist[v], nxt_bit_idx++);
//           count[dist[v]]++;
//           par[v] = -1;
//           dist[v] = -1;
//         }
//         bit_range.eb(cumsum<int>(count));
//         for (auto&& x: bit_range[K]) x += off;
//         st.eb(nbd, sz[nbd]);
//       }
//     }
//     // FenwickTree
//     vc<X> bit_raw(nxt_bit_idx);
//     FOR(v, N) {
//       for (auto&& [k, x, i]: dat[v]) { bit_raw[i] = v_vals[v]; }
//     }
//     bit.build(bit_raw);
//   }
// };

// TODO

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	cs := NewContourSum(nums)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		cs.AddEdge(u, v, w)
	}
	cs.Build()

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var root, add int
			fmt.Fscan(in, &root, &add)
			cs.Add(root, add)
		} else {
			var root, left, right int
			fmt.Fscan(in, &root, &left, &right)
			fmt.Fprintln(out, cs.Sum(root, left, right))
		}
	}

	// bit := NewBinaryIndexedTree(10)
	// bit.Apply(0, 1)
	// bit.Apply(1, 2)
	// fmt.Println(bit.ProdRange(0, 2))
}

type E = int

type ContourSum struct {
	n        int
	graph    [][][2]int
	values   []E
	bit      *BinaryIndexedTree
	bitRange [][]int    // centroid ごと、方向ごと
	dat      [][][3]int // 方向ラベル、重心からの距離、bit でのindex
}

func NewContourSum(values []E) *ContourSum {
	res := &ContourSum{}
	res.n = len(values)
	res.graph = make([][][2]int, res.n)
	res.values = values
	return res
}

func (cs *ContourSum) AddEdge(u, v, w int) {
	cs.graph[u] = append(cs.graph[u], [2]int{v, w})
	cs.graph[v] = append(cs.graph[v], [2]int{u, w})
}

func (cs *ContourSum) Build() {
	nextBitIndex := 0
	done := make([]bool, cs.n)
	sz := make([]int, cs.n)
	par, dist := make([]int, cs.n), make([]int, cs.n)
	for i := 0; i < cs.n; i++ {
		par[i], dist[i] = -1, -1
	}

	cs.bitRange = make([][]int, cs.n)
	cs.dat = make([][][3]int, cs.n)
	st := [][2]int{{0, cs.n}}

	for len(st) > 0 {
		v0 := st[len(st)-1][0]
		n := st[len(st)-1][1]
		st = st[:len(st)-1]
		c := -1
		{
			var dfs func(v int) int
			dfs = func(v int) int {
				sz[v] = 1
				for _, e := range cs.graph[v] {
					to := e[0]
					if to != par[v] && !done[to] {
						par[to] = v
						sz[v] += dfs(to)
					}
				}

				if c == -1 && n-sz[v] <= n/2 {
					c = v
				}
				return sz[v]
			}
			dfs(v0)
		}

		// center からの bfs。部分木サイズもとっておく。
		done[c] = true
		{
			off := nextBitIndex
			que := []int{}
			add := func(v, d, p int) {
				if dist[v] != -1 {
					return
				}
				sz[v] = 1
				dist[v] = d
				par[v] = p
				que = append(que, v)
			}
			p := 0
			add(c, 0, -1)
			for p < len(que) {
				v := que[p]
				p++
				for _, e := range cs.graph[v] {
					to := e[0]
					if done[to] {
						continue
					}
					add(to, dist[v]+1, v)
				}
			}
			for i := len(que) - 1; i >= 1; i-- {
				v := que[i]
				sz[par[v]] += sz[v]
			}

			// 距離ごとのカウント
			maxD := dist[que[len(que)-1]]
			count := make([]int, maxD+1)
			// 重心、方向ラベル、重心からの距離、bit でのindex
			for _, v := range que {
				cs.dat[v] = append(cs.dat[v], [3]int{c, dist[v], nextBitIndex})
				nextBitIndex++
				count[dist[v]]++
				par[v] = -1
				dist[v] = -1
			}
			cs.bitRange[c] = make([]int, maxD+1)
			for i := 0; i < maxD; i++ {
				cs.bitRange[c][i+1] = cs.bitRange[c][i] + count[i]
			}
			for i := 0; i < maxD; i++ {
				cs.bitRange[c][i] += off
			}
		}

		// 方向ごとの bfs
		for _, e := range cs.graph[c] {
			off := nextBitIndex
			nbd := e[0]
			if done[nbd] {
				continue
			}
			k := len(cs.bitRange) // !TODO
			que := []int{}
			add := func(v, d int) {
				if dist[v] != -1 || v == c {
					return
				}
				dist[v] = d
				que = append(que, v)
			}
			p := 0
			add(nbd, 0)
			for p < len(que) {
				v := que[p]
				p++
				for _, e := range cs.graph[v] {
					to := e[0]
					if done[to] {
						continue
					}
					add(to, dist[v]+1)
				}
			}

			// 距離ごとのカウント
			maxD := dist[que[len(que)-1]]
			count := make([]int, maxD+1)
			// 重心、方向ラベル、重心からの距離、bit でのindex
			for _, v := range que {
				cs.dat[v] = append(cs.dat[v], [3]int{^k, dist[v], nextBitIndex})
				nextBitIndex++
				count[dist[v]]++
				dist[v] = -1
				par[v] = -1
			}

			tmp := make([]int, maxD+1)
			for i := 0; i < maxD; i++ {
				tmp[i+1] = tmp[i] + count[i]
			}
			cs.bitRange = append(cs.bitRange, tmp)
			for i := range cs.bitRange[k] {
				cs.bitRange[k][i] += off
			}
			st = append(st, [2]int{nbd, sz[nbd]})
		}
	}

	bitRaw := make([]E, nextBitIndex)
	for i := 0; i < cs.n; i++ {
		for _, e := range cs.dat[i] {
			bitRaw[e[2]] = cs.values[i]
		}
	}

	cs.bit = NewBinaryIndexedTreeFrom(bitRaw)

}

func (cs *ContourSum) Add(root int, value E) {
	for _, e := range cs.dat[root] {
		cs.bit.Apply(e[2], value)
	}
}

// root を中心として、距離 [left, right) の範囲の和
func (cs *ContourSum) Sum(root, left, right int) E {
	sum := e()
	for _, e := range cs.dat[root] {
		k, x := e[0], e[1]
		lo, hi := left-x, right-x
		k_ := k
		if k < 0 {
			lo, hi, k_ = lo-2, hi-2, ^k
		}
		n := len(cs.bitRange[k_]) - 2
		lo = max(lo, 0)
		hi = min(hi, n+1)
		if lo >= hi {
			continue
		}

		a, b := cs.bitRange[k_][lo], cs.bitRange[k_][hi]
		val := cs.bit.ProdRange(a, b)
		if k < 0 {
			val = inv(val)
		}
		sum = op(sum, val)
	}

	return sum
}

func e() E        { return 0 }
func op(a, b E) E { return a + b }
func inv(e E) E   { return -e }

type BinaryIndexedTree struct {
	n    int
	log  int
	data []int
}

// 長さ n の 0で初期化された配列で構築する.
func NewBinaryIndexedTree(n int) *BinaryIndexedTree {
	return &BinaryIndexedTree{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

// 配列で構築する.
func NewBinaryIndexedTreeFrom(arr []E) *BinaryIndexedTree {
	res := NewBinaryIndexedTree(len(arr))
	res.Build(arr)
	return res
}

func (b *BinaryIndexedTree) Build(arr []E) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] = op(b.data[j], b.data[i])
		}
	}
}

// 要素 i に値 v を加える.
func (b *BinaryIndexedTree) Apply(i int, v E) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] = op(b.data[i], v)
	}
}

// [0, r) の要素の総和を求める.
func (b *BinaryIndexedTree) Prod(r int) E {
	res := e()
	for ; r > 0; r -= r & -r {
		res = op(res, b.data[r])
	}
	return res
}

// [l, r) の要素の総和を求める.
func (b *BinaryIndexedTree) ProdRange(l, r int) E {
	l = max(l, 0)
	r = min(r, b.n)
	if l == 0 {
		return b.Prod(r)
	}
	pos, neg := e(), e()
	for l < r {
		pos = op(pos, b.data[r-1])
		r -= r & -r
	}
	for r < l {
		neg = op(neg, b.data[l-1])
		l -= l & -l
	}
	return op(pos, inv(neg))
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
