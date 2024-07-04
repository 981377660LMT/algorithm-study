// IncrementalSCC

package main

import (
	"bufio"
	"fmt"
	"os"
)

// // https://codeforces.com/blog/entry/91608
// // https://codeforces.com/contest/1989/problem/F
// // グラフの辺番号 0, 1, 2, ... 順に辺を足していく.
// // 各辺 i に対してそれがサイクルに含まれるような時刻の最小値 or infty を返す.
// // これで mst を作って path max query すれば 2 点が同じ scc になる時刻も求まる
// template <typename GT>
// vc<int> incremental_scc(GT& G) {
//   static_assert(GT::is_directed);
//   int N = G.N, M = G.M;
//   vc<int> merge_time(M, infty<int>);
//   vc<tuple<int, int, int>> dat;
//   FOR(i, M) {
//     auto& e = G.edges[i];
//     dat.eb(i, e.frm, e.to);
//   }

//   vc<int> new_idx(N, -1);
//   // L 時点ではサイクルには含まれず, R 時点では含まれる
//   auto dfs
//       = [&](auto& dfs, vc<tuple<int, int, int>>& dat, int L, int R) -> void {
//     if (dat.empty() || R == L + 1) return;
//     int M = (L + R) / 2;
//     int n = 0;
//     for (auto& [i, a, b]: dat) {
//       if (new_idx[a] == -1) new_idx[a] = n++;
//       if (new_idx[b] == -1) new_idx[b] = n++;
//     }

//	    Graph<int, 1> G(n);
//	    for (auto& [i, a, b]: dat) {
//	      if (i < M) G.add(new_idx[a], new_idx[b]);
//	    }
//	    G.build();
//	    auto [nc, comp] = strongly_connected_component(G);
//	    vc<tuple<int, int, int>> dat1, dat2;
//	    for (auto [i, a, b]: dat) {
//	      a = new_idx[a], b = new_idx[b];
//	      if (i < M) {
//	        if (comp[a] == comp[b]) {
//	          chmin(merge_time[i], M), dat1.eb(i, a, b);
//	        } else {
//	          dat2.eb(i, comp[a], comp[b]);
//	        }
//	      } else {
//	        dat2.eb(i, comp[a], comp[b]);
//	      }
//	    }
//	    for (auto& [i, a, b]: dat) new_idx[a] = new_idx[b] = -1;
//	    dfs(dfs, dat1, L, M), dfs(dfs, dat2, M, R);
//	  };
//	  dfs(dfs, dat, 0, M + 1);
//	  return merge_time;
//	}
// void solve() {
//   INT(N, M);
//   VEC(mint, X, N);

//   Graph<int, 1> G(N);
//   G.read_graph(M, 0, 0);

//   auto time = incremental_scc(G);
//   vc<vc<int>> IDS(M + 1);
//   FOR(i, M) {
//     if (time[i] != infty<int>) IDS[time[i]].eb(i);
//   }

//	  UnionFind uf(N);
//	  mint ANS = 0;
//	  FOR(t, 1, M + 1) {
//	    for (auto &i: IDS[t]) {
//	      int a = G.edges[i].frm;
//	      int b = G.edges[i].to;
//	      a = uf[a], b = uf[b];
//	      if (a == b) continue;
//	      ANS += X[a] * X[b];
//	      uf.merge(a, b);
//	      X[uf[a]] = X[a] + X[b];
//	    }
//	    print(ANS);
//	  }
//	}
//

const MOD int = 998244353

// https://judge.yosupo.jp/problem/incremental_scc
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	graph := make([][]int32, n)
	edges := make([][2]int32, 0, m)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		edges = append(edges, [2]int32{u, v})
		graph[u] = append(graph[u], v)
	}

	time := IncrementalScc(graph, edges)
	groupByTime := make([][]int32, m+1)
	for i, t := range time {
		if t != INF32 {
			groupByTime[t] = append(groupByTime[t], int32(i))
		}
	}

	uf := NewUnionFindArraySimple32(n)
	res := 0
	for t := int32(1); t <= m; t++ {
		eids := groupByTime[t]
		for _, eid := range eids {
			e := edges[eid]
			a, b := uf.Find(e[0]), uf.Find(e[1])
			if a == b {
				continue
			}
			res += nums[a] * nums[b]
			uf.Union(a, b)
			nums[uf.Find(a)] = nums[a] + nums[b]
		}
	}
}

const INF32 int32 = 1e9 + 10

// 按照顺序连接编号为 0, 1, 2, ... 的边。
// 对于每条边 i，返回它在包含的环中的最小时间.
// 如果它不在环中，返回 INF32.
// !用途：mst + path max query 可以求出两点首次在同一个 scc 中的时刻.
func IncrementalScc(directedGraph [][]int32, directedEdges [][2]int32) (mergeTime []int32) {
	n, m := int32(len(directedGraph)), int32(len(directedEdges))
	mergeTime = make([]int32, m)
	for i := range mergeTime {
		mergeTime[i] = INF32
	}
	data := make([][3]int32, 0, m) // (i,from,to)
	for i, e := range directedEdges {
		data = append(data, [3]int32{int32(i), e[0], e[1]})
	}

	newId := make([]int32, n)
	for i := range newId {
		newId[i] = -1
	}

	// L 时刻不在环中，R 时刻在环中
	var dfs func([][3]int32, int32, int32)
	dfs = func(data [][3]int32, L, R int32) {
		if len(data) == 0 || R == L+1 {
			return
		}
		mid := (L + R) >> 1
		n := int32(0)
		for j := range data {
			a, b := data[j][1], data[j][2]
			if newId[a] == -1 {
				newId[a] = n
				n++
			}
			if newId[b] == -1 {
				newId[b] = n
				n++
			}
		}
		newGraph := make([][]int32, n)
		for j := range data {
			i, a, b := data[j][0], data[j][1], data[j][2]
			if i < mid {
				newGraph[newId[a]] = append(newGraph[newId[a]], newId[b])
			}
		}
		_, belong := StronglyConnectedComponent(newGraph)
		var dat1, dat2 [][3]int32
		for j := range data {
			i, a, b := data[j][0], data[j][1], data[j][2]
			a, b = newId[a], newId[b]
			if i < mid {
				if belong[a] == belong[b] {
					if mid < mergeTime[i] {
						mergeTime[i] = mid
					}
					dat1 = append(dat1, [3]int32{i, a, b})
				} else {
					dat2 = append(dat2, [3]int32{i, belong[a], belong[b]})
				}
			} else {
				dat2 = append(dat2, [3]int32{i, belong[a], belong[b]})
			}
		}
		for j := range data {
			a, b := data[j][1], data[j][2]
			newId[a], newId[b] = -1, -1
		}
		dfs(data, L, mid)
		dfs(data, mid, R)
	}

	dfs(data, 0, m+1)
	return
}

func StronglyConnectedComponent(directedGraph [][]int32) (count int32, belong []int32) {
	n := int32(len(directedGraph))
	comp, low, ord := make([]int32, n), make([]int32, n), make([]int32, n)
	for i := range ord {
		ord[i] = -1
	}
	path := make([]int32, 0)
	now := int32(0)
	var dfs func(int32)
	dfs = func(v int32) {
		ord[v] = now
		low[v] = now
		now++
		path = append(path, v)
		for _, to := range directedGraph[v] {
			if ord[to] == -1 {
				dfs(to)
				if low[to] < low[v] {
					low[v] = low[to]
				}
			} else {
				if ord[to] < low[v] {
					low[v] = ord[to]
				}
			}
		}
		if low[v] == ord[v] {
			for {
				cur := path[len(path)-1]
				path = path[:len(path)-1]
				ord[cur] = n
				comp[cur] = count
				if cur == v {
					break
				}
			}
			count++
		}
	}
	for i := int32(0); i < n; i++ {
		if ord[i] == -1 {
			dfs(i)
		}
	}
	for i := int32(0); i < n; i++ {
		comp[i] = count - 1 - comp[i]
	}
	return
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32, preMerge func(big, small int32)) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	if preMerge != nil {
		preMerge(root1, root2)
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = u.Find(u.data[key])
	return u.data[key]
}

func (u *UnionFindArraySimple32) GetSize(key int32) int32 {
	return -u.data[u.Find(key)]
}
