// https://maspypy.github.io/library/graph/blackbox/mst.hpp


#include "ds/unionfind/unionfind.hpp"

// Brouvka
// find_unused(v)：unused なうちで、v と最小コストで結べる点を探す。
// pair<int,COST> なければ {-1,*} を返すこと。
template <typename COST, typename F0, typename F1, typename F2>
vc<tuple<int, int, COST>> blackbox_mst(int N, F0 set_used, F1 set_unused,
                                       F2 find_unused) {
  using edge = tuple<int, int, COST>;
  UnionFind uf(N);
  vc<edge> res;
  while (1) {
    bool upd = 0;
    vvc<int> comp(N);
    vc<edge> cand(N, {-1, -1, -1});
    FOR(v, N) comp[uf[v]].eb(v);
    FOR(v, N) if (uf[v] == v) {
      for (auto&& x: comp[v]) { set_used(x); }
      for (auto&& x: comp[v]) {
        auto [y, cost] = find_unused(x);
        if (y == -1) continue;
        auto [a, b, c] = cand[v];
        if (a == -1 || cost < c) { cand[v] = {x, y, cost}; }
      }
      for (auto&& x: comp[v]) { set_unused(x); }
    }
    FOR(v, N) if (uf[v] == v) {
      auto [a, b, c] = cand[v];
      if (a == -1) continue;
      upd = 1;
      if (uf.merge(a, b)) res.eb(a, b, c);
    }
    if (!upd) break;
  }
  return res;
}

// Brouvka
// https://atcoder.jp/contests/pakencamp-2023-day1/tasks/pakencamp_2023_day1_p
// incremental なデータ構造を使う版
// init：探索のためのデータ構造初期化等. 2logN 回程度呼ばれる.
// add：探索対象に追加
// find：(w,wt) or (-1,*)
template <typename COST, typename F0, typename F1, typename F2>
vc<tuple<int, int, COST>> blackbox_mst_incremental(int N, F0 init, F1 add,
                                                   F2 find) {
  using edge = tuple<int, int, COST>;
  UnionFind uf(N);
  vc<edge> res;
  while (uf.n_comp > 1) {
    bool upd = 0;
    vvc<int> comp(N);
    FOR(v, N) comp[uf[v]].eb(v);
    vc<COST> wt(N, infty<COST>);
    vc<pair<int, int>> who(N, {-1, -1});
    init();
    FOR(i, N) {
      for (auto& v: comp[i]) {
        auto [w, x] = find(v);
        if (chmin(wt[i], x)) who[i] = {v, w};
      }
      for (auto& v: comp[i]) { add(v); }
    }
    init();
    FOR_R(i, N) {
      for (auto& v: comp[i]) {
        auto [w, x] = find(v);
        if (chmin(wt[i], x)) who[i] = {v, w};
      }
      for (auto& v: comp[i]) { add(v); }
    }
    FOR(i, N) {
      auto [a, b] = who[i];
      if (a == -1) continue;
      if (uf.merge(a, b)) {
        upd = true;
        res.eb(a, b, wt[i]);
      }
    }
    if (!upd) break;
  }
  return res;
}

// Prim
// O(N^2)
// https://atcoder.jp/contests/pakencamp-2023-day1/tasks/pakencamp_2023_day1_p
template <typename COST, typename F>
vc<tuple<int, int, COST>> mst_prim(int N, F cost) {
  using edge = tuple<int, int, COST>;
  UnionFind uf(N);
  vc<edge> res;
  vc<ll> wt(N, infty<ll>);
  vc<int> TO(N);
  auto add = [&](int v) -> void {
    FOR(w, N) {
      if (TO[w] != -1 && chmin(wt[w], cost(v, w))) TO[w] = v;
    }
    wt[v] = infty<ll>;
    TO[v] = -1;
  };
  add(0);
  FOR(N - 1) {
    int v = min_element(all(wt)) - wt.begin();
    res.eb(TO[v], v, wt[v]);
    add(v);
  }
  return res;
}