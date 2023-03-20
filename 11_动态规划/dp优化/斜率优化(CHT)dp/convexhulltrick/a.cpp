template <typename T>
struct Line {
  mutable T k, m, p;
  bool operator<(const Line& o) const { return k < o.k; }
  bool operator<(T x) const { return p < x; }
};

template <typename T>
T lc_inf() {
  return numeric_limits<T>::max();
}
template <>
long double lc_inf<long double>() {
  return 1 / .0;
}

template <typename T>
T lc_div(T a, T b) {
  return a / b - ((a ^ b) < 0 and a % b);
}
template <>
long double lc_div(long double a, long double b) {
  return a / b;
};
template <>
double lc_div(double a, double b) {
  return a / b;
};

template <typename T, bool MINIMIZE = true>
struct LineContainer : multiset<Line<T>, less<>> {
  using super = multiset<Line<T>, less<>>;
  using super::begin, super::end, super::insert, super::erase;
  using super::empty, super::lower_bound;
  const T inf = lc_inf<T>();
  bool insect(typename super::iterator x, typename super::iterator y) {
    if (y == end()) return x->p = inf, false;
    if (x->k == y->k)
      x->p = (x->m > y->m ? inf : -inf);
    else
      x->p = lc_div(y->m - x->m, x->k - y->k);
    return x->p >= y->p;
  }
  void add(T k, T m) {
    if (MINIMIZE) { k = -k, m = -m; }
    auto z = insert({k, m, 0}), y = z++, x = y;
    while (insect(y, z)) z = erase(z);
    if (x != begin() and insect(--x, y)) insect(x, y = erase(y));
    while ((y = x) != begin() and (--x)->p >= y->p) insect(x, erase(y));
  }
  T query(T x) {
    assert(!empty());
    auto l = *lower_bound(x);
    T v = (l.k * x + l.m);
    return (MINIMIZE ? -v : v);
  }
};

template <typename T>
using CHT_min = LineContainer<T, true>;
template <typename T>
using CHT_max = LineContainer<T, false>;

/*
long long / double で動くと思う。クエリあたり O(log N)
・add(a, b)：ax + by の追加
・get_max(x,y)：max_{a,b} (ax + by)
・get_min(x,y)：max_{a,b} (ax + by)
*/
template <typename T>
struct CHT_xy {
  using ld = long double;
  CHT_min<ld> cht_min;
  CHT_max<ld> cht_max;
  T amax = -infty<T>, amin = infty<T>;
  T bmax = -infty<T>, bmin = infty<T>;
  bool empty = true;

  void clear() {
    empty = true;
    cht_min.clear();
    cht_max.clear();
  }
  void add(T a, T b) {
    empty = false;
    cht_min.add(b, a);
    cht_max.add(b, a);
    chmax(amax, a), chmin(amin, a), chmax(bmax, b), chmin(bmin, b);
  }

  T get_max(T x, T y) {
    if (cht_min.empty()) return -infty<T>;
    if (x == 0) { return max(bmax * y, bmin * y); }
    ld z = ld(y) / x;
    if (x > 0) {
      auto l = cht_max.lower_bound(z);
      ll a = l->m, b = l->k;
      return a * x + b * y;
    }
    auto l = cht_min.lower_bound(z);
    ll a = -(l->m), b = -(l->k);
    return a * x + b * y;
  }

  T get_min(T x, T y) { return -get_max(-x, -y); }
};