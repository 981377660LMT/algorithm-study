#pragma once

template <typename E, typename Id, E (*f)(E, E), E (*mapping)(E, Id), Id (*composition)(Id, Id),
          E (*e)(), Id (*id)()>
struct DynamicLazySegmentTree {
  using ll = long long;
  struct Node;
  using Ptr = Node *;

  struct Node {
    Ptr l, r;
    E sum;
    Id laz;
    Node() : l(nullptr), r(nullptr), sum(e()), laz(id()) {}
  };

  Ptr root;
  ll Rmax;
  DynamicLazySegmentTree(ll n = 0) {
    Rmax = 2;
    while (Rmax <= n) Rmax *= 2;
    root = my_new();
  }

  // get a[x]
  E get_val(ll x) {
    ll L = 0, R = Rmax;
    Ptr t = root;
    while (L + 1 < R) {
      pushDown(t);
      ll M = (L + R) / 2;
      if (x < M) {
        if (!t->l) return e();
        t = t->l, R = M;
      } else {
        if (!t->r) return e();
        t = t->r, L = M;
      }
    }
    return t->sum;
  }
  // a[i] <- x
  void set_val(ll i, const E &x) { _set_val(root, 0, Rmax, i, x); }
  // apply x to a[l], a[l+1], ..., a[r-1]
  void apply(ll l, ll r, const Id &x) { _apply(root, 0, Rmax, l, r, x); }
  // get sum(a[l], a[l+1], ..., a[r-1])
  E fold(ll l, ll r) { return _fold(root, 0, Rmax, l, r); }
  // 破壊的。また、key が等しい要素がある場合は未定義
  void merge(DynamicLazySegmentTree &rhs) { root = _merge(root, rhs.root); }

 private:
  DynamicLazySegmentTree(Ptr _root, ll _Rmax) : root(_root), Rmax(_Rmax) {}

  Ptr my_new() { return new Node{}; }
  void my_del(Ptr p) { delete p; }

  void propagate(Ptr t, const Id &x) {
    assert(t != nullptr && x != id());
    t->laz = composition(t->laz, x);
    t->sum = mapping(t->sum, x);
  }
  void pushDown(Ptr t) {
    assert(t != nullptr);
    if (t->laz == id()) return;
    if (!t->l) t->l = my_new();
    if (!t->r) t->r = my_new();
    propagate(t->l, t->laz);
    propagate(t->r, t->laz);
    t->laz = id();
  }
  void pushUp(Ptr t) {
    assert(t->laz == id());
    t->sum = f(t->l ? t->l->sum : e(), t->r ? t->r->sum : e());
  }

  void _set_val(Ptr t, ll L, ll R, ll i, const E &x) {
    assert(L <= i && i < R && t);
    if (L + 1 == R) {
      t->sum = x;
      return;
    }
    ll M = (L + R) / 2;
    pushDown(t);
    if (i < M) {
      if (!t->l) t->l = my_new();
      _set_val(t->l, L, M, i, x);
    } else {
      if (!t->r) t->r = my_new();
      _set_val(t->r, M, R, i, x);
    }
    pushUp(t);
  }

  void _apply(Ptr t, ll L, ll R, ll a, ll b, const Id &x) {
    assert(a <= b && a < R && L < b && t);
    if (L == a and R == b) {
      propagate(t, x);
      if (L + 1 == R) t->laz = id();
      return;
    }
    ll M = (L + R) / 2;
    pushDown(t);
    if (a < M) {
      if (!t->l) t->l = my_new();
      _apply(t->l, L, M, a, min(b, M), x);
    }
    if (M < b) {
      if (!t->r) t->r = my_new();
      _apply(t->r, M, R, max(a, M), b, x);
    }
    pushUp(t);
    return;
  }

  E _fold(Ptr t, ll L, ll R, ll a, ll b) {
    assert(a <= b && a < R && L < b && t);
    if (L == a and R == b) return t->sum;
    ll M = (L + R) / 2;
    pushDown(t);
    E v = e();
    if (a < M && t->l) v = f(_fold(t->l, L, M, a, min(b, M)), v);
    if (M < b && t->r) v = f(v, _fold(t->r, M, R, max(a, M), b));
    return v;
  }

  Ptr _merge(Ptr t1, Ptr t2) {
    if (!t1 or !t2) return t1 ? t1 : t2;
    assert(t1->laz == id() && t2->laz == id());
    t1->l = merge(t1->l, t2->l);
    t1->r = merge(t1->r, t2->r);
    pushUp(t1), my_del(t2);
    return t1;
  }

  /*
  // [L, x), [x, R) で split
  pair<Ptr, Ptr> _split(Ptr t1, ll L, ll R, ll x) {
    if (!t1) return {nullptr, nullptr};
    assert(t1->laz == id());
    Ptr t2 = my_new();
    ll M = (L + R) / 2;
  }
  */
};

namespace DynamicSegmentTreeImpl {
template <typename E>
E mapping(E l, bool) {
  return l;
}
bool composition(bool, bool) { return false; }
bool id() { return false; }

template <typename E, E (*f)(E, E), E (*e)()>
using DynamicSegmentTree = DynamicLazySegmentTree<E, bool, f, mapping, composition, e, id>;
}  // namespace DynamicSegmentTreeImpl

using DynamicSegmentTreeImpl::DynamicSegmentTree;

/**
 * @brief Dynamic Segment Tree
 */