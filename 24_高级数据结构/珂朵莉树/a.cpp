
// https://codeforces.com/contest/1638/problem/E
// 持つ値のタイプ T、座標タイプ X
// コンストラクタでは T none_val を指定する
template <typename T, typename X = ll>
struct Intervals {
  static constexpr X LLIM = -infty<X>;
  static constexpr X RLIM = infty<X>;
  const T none_val;
  // none_val でない区間の個数と長さ合計
  int total_num;
  X total_len;
  map<X, T> dat;

  Intervals(T none_val) : none_val(none_val), total_num(0), total_len(0) {
    dat[LLIM] = none_val;
    dat[RLIM] = none_val;
  }

  // x を含む区間の情報の取得
  tuple<X, X, T> get(X x, bool ERASE) {
    auto it2 = dat.upper_bound(x);
    auto it1 = prev(it2);
    auto [l, tl] = *it1;
    auto [r, tr] = *it2;
    if (tl != none_val && ERASE) {
      --total_num, total_len -= r - l;
      dat[l] = none_val;
      merge_at(l);
      merge_at(r);
    }
    return {l, r, tl};
  }

  // [L, R) 内の全データの取得
  template <typename F>
  void enumerate_range(X L, X R, F f, bool ERASE) {
    assert(LLIM <= L && L <= R && R <= RLIM);
    if (!ERASE) {
      auto it = prev(dat.upper_bound(L));
      while ((*it).fi < R) {
        auto it2 = next(it);
        f(max((*it).fi, L), min((*it2).fi, R), (*it).se);
        it = it2;
      }
      return;
    }
    // 半端なところの分割
    auto p = prev(dat.upper_bound(L));
    if ((*p).fi < L) {
      dat[L] = (*p).se;
      if (dat[L] != none_val) ++total_num;
    }
    p = dat.lower_bound(R);
    if (R < (*p).fi) {
      T t = (*prev(p)).se;
      dat[R] = t;
      if (t != none_val) ++total_num;
    }
    p = dat.lower_bound(L);
    while (1) {
      if ((*p).fi >= R) break;
      auto q = next(p);
      T t = (*p).se;
      f((*p).fi, (*q).fi, t);
      if (t != none_val) --total_num, total_len -= (*q).fi - (*p).fi;
      p = dat.erase(p);
    }
    dat[L] = none_val;
  }

  void set(X L, X R, T t) {
    enumerate_range(
        L, R, [](int l, int r, T x) -> void {}, true);
    dat[L] = t;
    if (t != none_val) total_num++, total_len += R - L;
    merge_at(L);
    merge_at(R);
  }

  template <typename F>
  void enumerate_all(F f) {
    enumerate_range(LLIM, RLIM, f, false);
  }

  void merge_at(X p) {
    if (p == LLIM || RLIM == p) return;
    auto itp = dat.lower_bound(p);
    assert((*itp).fi == p);
    auto itq = prev(itp);
    if ((*itp).se == (*itq).se) {
      if ((*itp).se != none_val) --total_num;
      dat.erase(itp);
    }
  }
};