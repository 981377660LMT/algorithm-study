#include "graph/base.hpp"
#include "enumerate/bits.hpp"

// 辺重みは e.cost、頂点重みは vector で渡す。返り値：{cost, vs, es}
// O(3^kn + 2^k(n+m)log n), k: terminal size
template <typename T, typename GT>
tuple<T, vc<int>, vc<int>> steiner_tree(GT& G, vc<int> S, vc<T> v_wt = {}) {
  assert(!S.empty() && !G.is_directed());
  const int N = G.N, M = G.M, K = len(S);
  if (v_wt.empty()) v_wt.assign(N, 0);

  // ターミナル集合, root -> cost
  vv(T, DP, 1 << K, N, infty<T>);
  FOR(v, N) DP[0][v] = v_wt[v];

  // 2 * t or 2 * eid + 1
  vv(int, par, 1 << K, N, -1);

  FOR(s, 1, 1 << K) {
    auto& dp = DP[s];
    enumerate_bits(s, [&](int k) -> void {
      int v = S[k];
      chmin(dp[v], DP[s ^ 1 << k][v]);
    });
    FOR_subset(t, s) {
      if (t == 0 || t == s) continue;
      FOR(v, N) {
        if (chmin(dp[v], DP[t][v] + DP[s ^ t][v] - v_wt[v])) par[s][v] = 2 * t;
      }
    }
    // 根の移動を dijkstra で
    pqg<pair<T, int>> que;
    FOR(v, N) que.emplace(dp[v], v);
    while (!que.empty()) {
      auto [dv, v] = POP(que);
      if (dv != dp[v]) continue;
      for (auto&& e: G[v]) {
        if (chmin(dp[e.to], dv + e.cost + v_wt[e.to])) {
          par[s][e.to] = 2 * e.id + 1;
          que.emplace(dp[e.to], e.to);
        }
      }
    }
  }

  // 復元する
  vc<bool> used_v(N), used_e(M);
  vc<int> v_to_k(N, -1);
  FOR(k, K) v_to_k[S[k]] = k;

  vc<pair<int, int>> que;
  int root = min_element(all(DP.back())) - DP.back().begin();
  que.eb((1 << K) - 1, root);
  used_v[root] = 1;

  while (len(que)) {
    auto [s, v] = POP(que);
    if (s == 0) { continue; }
    if (par[s][v] == -1) {
      int k = v_to_k[v];
      assert(k != -1 && s >> k & 1);
      que.eb(s ^ 1 << k, v);
      continue;
    }
    elif (par[s][v] & 1) {
      int eid = par[s][v] / 2;
      auto& e = G.edges[eid];
      int w = v ^ e.frm ^ e.to;
      used_v[w] = used_e[eid] = 1;
      que.eb(s, w);
      continue;
    }
    else {
      int t = par[s][v] / 2;
      que.eb(t, v), que.eb(s ^ t, v);
    }
  }
  vc<int> vs, es;
  FOR(v, N) if (used_v[v]) vs.eb(v);
  FOR(e, M) if (used_e[e]) es.eb(e);
  T cost = 0;
  for (auto&& v: vs) cost += v_wt[v];
  for (auto&& e: es) cost += G.edges[e].cost;
  assert(cost == DP.back()[root]);
  return {cost, vs, es};
}