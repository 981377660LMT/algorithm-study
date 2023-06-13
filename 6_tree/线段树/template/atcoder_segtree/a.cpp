#pragma once

template <typename T, typename F>
struct SegmentTree {
  int N;
  int size;
  vector<T> seg;
  const F f;
  const T I;

  SegmentTree(F _f, const T &I_) : N(0), size(0), f(_f), I(I_) {}

  SegmentTree(int _N, F _f, const T &I_) : f(_f), I(I_) { init(_N); }

  SegmentTree(const vector<T> &v, F _f, T I_) : f(_f), I(I_) {
    init(v.size());
    for (int i = 0; i < (int)v.size(); i++) {
      seg[i + size] = v[i];
    }
    build();
  }

  void init(int _N) {
    N = _N;
    size = 1;
    while (size < N) size <<= 1;
    seg.assign(2 * size, I);
  }

  void set(int k, T x) { seg[k + size] = x; }

  void build() {
    for (int k = size - 1; k > 0; k--) {
      seg[k] = f(seg[2 * k], seg[2 * k + 1]);
    }
  }

  void update(int k, T x) {
    k += size;
    seg[k] = x;
    while (k >>= 1) {
      seg[k] = f(seg[2 * k], seg[2 * k + 1]);
    }
  }

  void add(int k, T x) {
    k += size;
    seg[k] += x;
    while (k >>= 1) {
      seg[k] = f(seg[2 * k], seg[2 * k + 1]);
    }
  }

  // query to [a, b)
  T query(int a, int b) {
    T L = I, R = I;
    for (a += size, b += size; a < b; a >>= 1, b >>= 1) {
      if (a & 1) L = f(L, seg[a++]);
      if (b & 1) R = f(seg[--b], R);
    }
    return f(L, R);
  }

  T &operator[](const int &k) { return seg[k + size]; }

  // check(a[l] * ...  * a[r-1]) が true となる最大の r
  // (右端まですべて true なら N を返す)
  template <class C>
  int max_right(int l, C check) {
    assert(0 <= l && l <= N);
    assert(check(I) == true);
    if (l == N) return N;
    l += size;
    T sm = I;
    do {
      while (l % 2 == 0) l >>= 1;
      if (!check(f(sm, seg[l]))) {
        while (l < size) {
          l = (2 * l);
          if (check(f(sm, seg[l]))) {
            sm = f(sm, seg[l]);
            l++;
          }
        }
        return l - size;
      }
      sm = f(sm, seg[l]);
      l++;
    } while ((l & -l) != l);
    return N;
  }

  // check(a[l] * ... * a[r-1]) が true となる最小の l
  // (左端まで true なら 0 を返す)
  template <typename C>
  int min_left(int r, C check) {
    assert(0 <= r && r <= N);
    assert(check(I) == true);
    if (r == 0) return 0;
    r += size;
    T sm = I;
    do {
      r--;
      while (r > 1 && (r % 2)) r >>= 1;
      if (!check(f(seg[r], sm))) {
        while (r < size) {
          r = (2 * r + 1);
          if (check(f(seg[r], sm))) {
            sm = f(seg[r], sm);
            r--;
          }
        }
        return r + 1 - size;
      }
      sm = f(seg[r], sm);
    } while ((r & -r) != r);
    return 0;
  }
};