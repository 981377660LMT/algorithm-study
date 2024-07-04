// https://codeforces.com/blog/entry/91608
// https://codeforces.com/contest/1989/problem/F
// グラフの辺番号 0, 1, 2, ... 順に辺を足していく.
// 各辺 i に対してそれがサイクルに含まれるような時刻の最小値 or infty を返す.
// これで mst を作って path max query すれば 2 点が同じ scc になる時刻も求まる
template <typename GT>
vc<int> incremental_scc(GT& G) {
  static_assert(GT::is_directed);
  int N = G.N, M = G.M;
  vc<int> merge_time(M, infty<int>);
  vc<tuple<int, int, int>> dat;
  FOR(i, M) {
    auto& e = G.edges[i];
    dat.eb(i, e.frm, e.to);
  }

  vc<int> new_idx(N, -1);
  // L 時点ではサイクルには含まれず, R 時点では含まれる
  auto dfs
      = [&](auto& dfs, vc<tuple<int, int, int>>& dat, int L, int R) -> void {
    if (dat.empty() || R == L + 1) return;
    int M = (L + R) / 2;
    int n = 0;
    for (auto& [i, a, b]: dat) {
      if (new_idx[a] == -1) new_idx[a] = n++;
      if (new_idx[b] == -1) new_idx[b] = n++;
    }

    Graph<int, 1> G(n);
    for (auto& [i, a, b]: dat) {
      if (i < M) G.add(new_idx[a], new_idx[b]);
    }
    G.build();
    auto [nc, comp] = strongly_connected_component(G);
    vc<tuple<int, int, int>> dat1, dat2;
    for (auto [i, a, b]: dat) {
      a = new_idx[a], b = new_idx[b];
      if (i < M) {
        if (comp[a] == comp[b]) {
          chmin(merge_time[i], M), dat1.eb(i, a, b);
        } else {
          dat2.eb(i, comp[a], comp[b]);
        }
      } else {
        dat2.eb(i, comp[a], comp[b]);
      }
    }
    for (auto& [i, a, b]: dat) new_idx[a] = new_idx[b] = -1;
    dfs(dfs, dat1, L, M), dfs(dfs, dat2, M, R);
  };
  dfs(dfs, dat, 0, M + 1);
  return merge_time;
}