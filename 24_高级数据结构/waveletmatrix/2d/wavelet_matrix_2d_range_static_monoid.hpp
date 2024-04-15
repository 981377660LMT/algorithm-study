#include "ds/bit_vector.hpp"
#include "ds/segtree/segtree.hpp"
#include "alg/monoid/add.hpp"
#include "ds/static_range_product.hpp"

template <typename Monoid, typename ST, typename XY, bool SMALL_X, bool SMALL_Y>
struct Wavelet_Matrix_2D_Range_Static_Monoid {
  // 点群を Y 昇順に並べる.
  // X を整数になおして binary trie みたいに振り分ける
  using MX = Monoid;
  using X = typename MX::value_type;
  static_assert(MX::commute);

  template <bool SMALL>
  struct TO_IDX {
    vc<XY> key;
    XY mi, ma;
    vc<int> dat;

    void build(vc<XY>& X) {
      if constexpr (SMALL) {
        mi = (X.empty() ? 0 : MIN(X));
        ma = (X.empty() ? 0 : MAX(X));
        dat.assign(ma - mi + 2, 0);
        for (auto& x: X) { dat[x - mi + 1]++; }
        FOR(i, len(dat) - 1) dat[i + 1] += dat[i];
      } else {
        key = X;
        sort(all(key));
      }
    }
    int operator()(XY x) {
      if constexpr (SMALL) {
        return dat[clamp<XY>(x - mi, 0, ma - mi + 1)];
      } else {
        return LB(key, x);
      }
    }
  };

  TO_IDX<SMALL_X> XtoI;
  TO_IDX<SMALL_Y> YtoI;

  int N, lg;
  vector<int> mid;
  vector<Bit_Vector> bv;
  vc<int> new_idx;
  vc<int> A;
  using SEG = Static_Range_Product<Monoid, ST, 4>;
  vc<SEG> dat;

  template <typename F>
  Wavelet_Matrix_2D_Range_Static_Monoid(int N, F f) {
    build(N, f);
  }

  template <typename F>
  void build(int N_, F f) {
    N = N_;
    if (N == 0) {
      lg = 0;
      return;
    }
    vc<XY> tmp(N), Y(N);
    vc<X> S(N);
    FOR(i, N) tie(tmp[i], Y[i], S[i]) = f(i);
    auto I = argsort(Y);
    tmp = rearrange(tmp, I), Y = rearrange(Y, I), S = rearrange(S, I);
    XtoI.build(tmp), YtoI.build(Y);
    new_idx.resize(N);
    FOR(i, N) new_idx[I[i]] = i;

    // あとは普通に
    lg = tmp.empty() ? 0 : __lg(XtoI(MAX(tmp) + 1)) + 1;
    mid.resize(lg), bv.assign(lg, Bit_Vector(N));
    dat.resize(lg);
    A.resize(N);
    FOR(i, N) A[i] = XtoI(tmp[i]);

    vc<XY> A0(N), A1(N);
    vc<X> S0(N), S1(N);
    FOR_R(d, lg) {
      int p0 = 0, p1 = 0;
      FOR(i, N) {
        bool f = (A[i] >> d & 1);
        if (!f) { S0[p0] = S[i], A0[p0] = A[i], p0++; }
        if (f) { S1[p1] = S[i], A1[p1] = A[i], bv[d].set(i), p1++; }
      }
      mid[d] = p0;
      bv[d].build();
      swap(A, A0), swap(S, S0);
      FOR(i, p1) A[p0 + i] = A1[i], S[p0 + i] = S1[i];
      dat[d].build(N, [&](int i) -> X { return S[i]; });
    }
    FOR(i, N) A[i] = XtoI(tmp[i]);
  }

  int count(XY x1, XY x2, XY y1, XY y2) {
    if (N == 0) return 0;
    x1 = XtoI(x1), x2 = XtoI(x2);
    y1 = YtoI(y1), y2 = YtoI(y2);
    return prefix_count(y1, y2, x2) - prefix_count(y1, y2, x1);
  }

  X prod(XY x1, XY x2, XY y1, XY y2) {
    if (N == 0) return MX::unit();
    assert(x1 <= x2 && y1 <= y2);
    x1 = XtoI(x1), x2 = XtoI(x2);
    y1 = YtoI(y1), y2 = YtoI(y2);
    X res = MX::unit();
    prod_dfs(y1, y2, x1, x2, lg - 1, res);
    return res;
  }

private:
  int prefix_count(int L, int R, int x) {
    int cnt = 0;
    FOR_R(d, lg) {
      int l0 = bv[d].rank(L, 0), r0 = bv[d].rank(R, 0);
      if (x >> d & 1) {
        cnt += r0 - l0, L += mid[d] - l0, R += mid[d] - r0;
      } else {
        L = l0, R = r0;
      }
    }
    return cnt;
  }

  void prod_dfs(int L, int R, int x1, int x2, int d, X& res) {
    chmax(x1, 0), chmin(x2, 1 << (d + 1));
    if (x1 >= x2) { return; }
    assert(0 <= x1 && x1 < x2 && x2 <= (1 << (d + 1)));
    if (x1 == 0 && x2 == (1 << (d + 1))) {
      res = MX::op(res, dat[d + 1].prod(L, R));
      return;
    }
    int l0 = bv[d].rank(L, 0), r0 = bv[d].rank(R, 0);
    prod_dfs(l0, r0, x1, x2, d - 1, res);
    prod_dfs(L + mid[d] - l0, R + mid[d] - r0, x1 - (1 << d), x2 - (1 << d),
             d - 1, res);
  }
};