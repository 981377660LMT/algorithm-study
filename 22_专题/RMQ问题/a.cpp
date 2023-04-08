// set,prod どちらも O(sqrt(N)) 時間。
// モノイドが可換なら普通のセグ木を使うこと。
template <class Monoid>
struct Xor_SegTree {
  static_assert(!Monoid::commute);
  using MX = Monoid;
  using X = typename MX::value_type;
  using value_type = X;
  vvc<X> dat;
  int n, log, size;
  int H; // 幅 2^H のところまで作る

  Xor_SegTree() {}
  Xor_SegTree(int n) { build(n); }
  template <typename F>
  Xor_SegTree(int n, F f) {
    build(n, f);
  }
  Xor_SegTree(const vc<X>& v) { build(v); }

  void build(int m) {
    build(m, [](int i) -> X { return MX::unit(); });
  }
  void build(const vc<X>& v) {
    build(len(v), [&](int i) -> X { return v[i]; });
  }
  template <typename F>
  void build(int m, F f) {
    n = m, log = 1;
    while ((1 << log) < n) ++log;
    size = 1 << log;
    assert(n == size);
    H = log / 2;
    dat.assign(H + 1, vc<X>(size));
    FOR(i, n) dat[0][i] = f(i);
    FOR(h, 1, H + 1) FOR(i, n >> h) { update(h, i); }
  }

  X get(int i) { return dat[0][i]; }
  vc<X> get_all() { return dat[0]; }

  void update(int h, int i) {
    assert(0 < h && h <= H);
    int cnt = 1 << (h - 1);
    int a = i << h;
    int b = a + cnt;
    FOR(k, cnt) {
      dat[h][a + k] = MX::op(dat[h - 1][a + k], dat[h - 1][b + k]);
      dat[h][b + k] = MX::op(dat[h - 1][b + k], dat[h - 1][a + k]);
    }
  }

  void set(int i, const X& x) {
    assert(i < n);
    dat[0][i] = x;
    FOR(h, 1, H + 1) {
      i /= 2;
      update(h, i);
    }
  }

  void multiply(int i, const X& x) {
    assert(i < n);
    dat[0][i] = MX::op(dat[0][i], x);
    FOR(h, 1, H + 1) {
      i /= 2;
      update(h, i);
    }
  }

  X prod(int L, int R, int xor_val) {
    X x1 = MX::unit(), x2 = MX::unit();
    FOR(h, H) {
      if (L >= R) break;
      if (L & (1 << h)) {
        x1 = MX::op(x1, dat[h][L ^ xor_val]);
        L += 1 << h;
      }
      if (R & (1 << h)) {
        R -= 1 << h;
        x2 = MX::op(dat[h][R ^ xor_val], x2);
      }
    }
    while (L < R) {
      x1 = MX::op(x1, dat[H][L ^ xor_val]);
      L += 1 << H;
    }
    return MX::op(x1, x2);
  }
};