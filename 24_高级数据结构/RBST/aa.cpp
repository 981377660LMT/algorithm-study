template <typename ActedMonoid, bool PERSISTENT, int NODES>
struct RBST_ActedMonoid {
  using Monoid_X = typename ActedMonoid::Monoid_X;
  using Monoid_A = typename ActedMonoid::Monoid_A;
  using X = typename Monoid_X::value_type;
  using A = typename Monoid_A::value_type;

  struct Node {
    Node *l, *r;
    X x, prod; // lazy, rev 反映済
    A lazy;
    u32 size;
    bool rev;
  };

  Node *pool;
  int pid;
  using np = Node *;

  RBST_ActedMonoid() : pid(0) { pool = new Node[NODES]; }

  void reset() { pid = 0; }

  np new_node(const X &x) {
    pool[pid].l = pool[pid].r = nullptr;
    pool[pid].x = x;
    pool[pid].prod = x;
    pool[pid].lazy = Monoid_A::unit();
    pool[pid].size = 1;
    pool[pid].rev = 0;
    return &(pool[pid++]);
  }

  np new_node(const vc<X> &dat) {
    auto dfs = [&](auto &dfs, u32 l, u32 r) -> np {
      if (l == r) return nullptr;
      if (r == l + 1) return new_node(dat[l]);
      u32 m = (l + r) / 2;
      np l_root = dfs(dfs, l, m);
      np r_root = dfs(dfs, m + 1, r);
      np root = new_node(dat[m]);
      root->l = l_root, root->r = r_root;
      update(root);
      return root;
    };
    return dfs(dfs, 0, len(dat));
  }

  np copy_node(np &n) {
    if (!n || !PERSISTENT) return n;
    pool[pid].l = n->l, pool[pid].r = n->r;
    pool[pid].x = n->x;
    pool[pid].prod = n->prod;
    pool[pid].lazy = n->lazy;
    pool[pid].size = n->size;
    pool[pid].rev = n->rev;
    return &(pool[pid++]);
  }

  np merge(np l_root, np r_root) { return merge_rec(l_root, r_root); }
  np merge3(np a, np b, np c) { return merge(merge(a, b), c); }
  np merge4(np a, np b, np c, np d) { return merge(merge(merge(a, b), c), d); }
  pair<np, np> split(np root, u32 k) {
    if (!root) {
      assert(k == 0);
      return {nullptr, nullptr};
    }
    assert(0 <= k && k <= root->size);
    return split_rec(root, k);
  }
  tuple<np, np, np> split3(np root, u32 l, u32 r) {
    np nm, nr;
    tie(root, nr) = split(root, r);
    tie(root, nm) = split(root, l);
    return {root, nm, nr};
  }
  tuple<np, np, np, np> split4(np root, u32 i, u32 j, u32 k) {
    np d;
    tie(root, d) = split(root, k);
    auto [a, b, c] = split3(root, i, j);
    return {a, b, c, d};
  }

  X prod(np root, u32 l, u32 r) {
    if (l == r) return Monoid_X::unit();
    return prod_rec(root, l, r, false);
  }
  X prod(np root) { return (root ? root->prod : Monoid_X::unit()); }

  np reverse(np root, u32 l, u32 r) {
    assert(Monoid_X::commute);
    assert(0 <= l && l <= r && r <= root->size);
    if (r - l <= 1) return root;
    auto [nl, nm, nr] = split3(root, l, r);
    nm->rev ^= 1;
    swap(nm->l, nm->r);
    return merge3(nl, nm, nr);
  }

  np apply(np root, u32 l, u32 r, const A a) {
    assert(0 <= l && l <= r && r <= root->size);
    return apply_rec(root, l, r, a);
  }
  np apply(np root, const A a) {
    if (!root) return root;
    return apply_rec(root, 0, root->size, a);
  }

  np set(np root, u32 k, const X &x) { return set_rec(root, k, x); }
  np multiply(np root, u32 k, const X &x) { return multiply_rec(root, k, x); }
  X get(np root, u32 k) { return get_rec(root, k, false, Monoid_A::unit()); }

  vc<X> get_all(np root) {
    vc<X> res;
    auto dfs = [&](auto &dfs, np root, bool rev, A lazy) -> void {
      if (!root) return;
      X me = ActedMonoid::act(root->x, lazy, 1);
      lazy = Monoid_A::op(root->lazy, lazy);
      dfs(dfs, (rev ? root->r : root->l), rev ^ root->rev, lazy);
      res.eb(me);
      dfs(dfs, (rev ? root->l : root->r), rev ^ root->rev, lazy);
    };
    dfs(dfs, root, 0, Monoid_A::unit());
    return res;
  }

  template <typename F>
  pair<np, np> split_max_right(np root, const F check) {
    assert(check(Monoid_X::unit()));
    X x = Monoid_X::unit();
    return split_max_right_rec(root, check, x);
  }

private:
  inline u32 xor128() {
    static u32 x = 123456789;
    static u32 y = 362436069;
    static u32 z = 521288629;
    static u32 w = 88675123;
    u32 t = x ^ (x << 11);
    x = y;
    y = z;
    z = w;
    return w = (w ^ (w >> 19)) ^ (t ^ (t >> 8));
  }

  void prop(np c) {
    // 自身をコピーする必要はない。
    // 子をコピーする必要がある。複数の親を持つ可能性があるため。
    bool bl_lazy = (c->lazy != Monoid_A::unit());
    bool bl_rev = c->rev;
    if (bl_lazy || bl_rev) {
      c->l = copy_node(c->l);
      c->r = copy_node(c->r);
    }
    if (c->lazy != Monoid_A::unit()) {
      if (c->l) {
        c->l->x = ActedMonoid::act(c->l->x, c->lazy, 1);
        c->l->prod = ActedMonoid::act(c->l->prod, c->lazy, c->l->size);
        c->l->lazy = Monoid_A::op(c->l->lazy, c->lazy);
      }
      if (c->r) {
        c->r->x = ActedMonoid::act(c->r->x, c->lazy, 1);
        c->r->prod = ActedMonoid::act(c->r->prod, c->lazy, c->r->size);
        c->r->lazy = Monoid_A::op(c->r->lazy, c->lazy);
      }
      c->lazy = Monoid_A::unit();
    }
    if (c->rev) {
      if (c->l) {
        c->l->rev ^= 1;
        swap(c->l->l, c->l->r);
      }
      if (c->r) {
        c->r->rev ^= 1;
        swap(c->r->l, c->r->r);
      }
      c->rev = 0;
    }
  }

  void update(np c) {
    // データを保ったまま正常化するだけなので、コピー不要
    c->size = 1;
    c->prod = c->x;
    if (c->l) {
      c->size += c->l->size;
      c->prod = Monoid_X::op(c->l->prod, c->prod);
    }
    if (c->r) {
      c->size += c->r->size;
      c->prod = Monoid_X::op(c->prod, c->r->prod);
    }
  }

  np merge_rec(np l_root, np r_root) {
    if (!l_root) return r_root;
    if (!r_root) return l_root;
    u32 sl = l_root->size, sr = r_root->size;
    if (xor128() % (sl + sr) < sl) {
      prop(l_root);
      l_root = copy_node(l_root);
      l_root->r = merge_rec(l_root->r, r_root);
      update(l_root);
      return l_root;
    }
    prop(r_root);
    r_root = copy_node(r_root);
    r_root->l = merge_rec(l_root, r_root->l);
    update(r_root);
    return r_root;
  }

  pair<np, np> split_rec(np root, u32 k) {
    if (!root) return {nullptr, nullptr};
    prop(root);
    u32 sl = (root->l ? root->l->size : 0);
    if (k <= sl) {
      auto [nl, nr] = split_rec(root->l, k);
      root = copy_node(root);
      root->l = nr;
      update(root);
      return {nl, root};
    }
    auto [nl, nr] = split_rec(root->r, k - (1 + sl));
    root = copy_node(root);
    root->r = nl;
    update(root);
    return {root, nr};
  }

  np set_rec(np root, u32 k, const X &x) {
    if (!root) return root;
    prop(root);
    u32 sl = (root->l ? root->l->size : 0);
    if (k < sl) {
      root = copy_node(root);
      root->l = set_rec(root->l, k, x);
      update(root);
      return root;
    }
    if (k == sl) {
      root = copy_node(root);
      root->x = x;
      update(root);
      return root;
    }
    root = copy_node(root);
    root->r = set_rec(root->r, k - (1 + sl), x);
    update(root);
    return root;
  }

  np multiply_rec(np root, u32 k, const X &x) {
    if (!root) return root;
    prop(root);
    u32 sl = (root->l ? root->l->size : 0);
    if (k < sl) {
      root = copy_node(root);
      root->l = multiply_rec(root->l, k, x);
      update(root);
      return root;
    }
    if (k == sl) {
      root = copy_node(root);
      root->x = Monoid_X::op(root->x, x);
      update(root);
      return root;
    }
    root = copy_node(root);
    root->r = multiply_rec(root->r, k - (1 + sl), x);
    update(root);
    return root;
  }

  X prod_rec(np root, u32 l, u32 r, bool rev) {
    if (l == 0 && r == root->size) { return root->prod; }
    np left = (rev ? root->r : root->l);
    np right = (rev ? root->l : root->r);
    u32 sl = (left ? left->size : 0);
    X res = Monoid_X::unit();
    if (l < sl) {
      X y = prod_rec(left, l, min(r, sl), rev ^ root->rev);
      res = Monoid_X::op(res, ActedMonoid::act(y, root->lazy, min(r, sl) - l));
    }
    if (l <= sl && sl < r) res = Monoid_X::op(res, root->x);
    u32 k = 1 + sl;
    if (k < r) {
      X y = prod_rec(right, max(k, l) - k, r - k, rev ^ root->rev);
      res = Monoid_X::op(res, ActedMonoid::act(y, root->lazy, r - max(k, l)));
    }
    return res;
  }

  X get_rec(np root, u32 k, bool rev, A lazy) {
    np left = (rev ? root->r : root->l);
    np right = (rev ? root->l : root->r);
    u32 sl = (left ? left->size : 0);
    if (k == sl) return ActedMonoid::act(root->x, lazy, 1);
    lazy = Monoid_A::op(root->lazy, lazy);
    rev ^= root->rev;
    if (k < sl) return get_rec(left, k, rev, lazy);
    return get_rec(right, k - (1 + sl), rev, lazy);
  }

  np apply_rec(np root, u32 l, u32 r, const A &a) {
    prop(root);
    root = copy_node(root);
    if (l == 0 && r == root->size) {
      root->x = ActedMonoid::act(root->x, a, 1);
      root->prod = ActedMonoid::act(root->prod, a, root->size);
      root->lazy = a;
      return root;
    }
    u32 sl = (root->l ? root->l->size : 0);
    if (l < sl) root->l = apply_rec(root->l, l, min(r, sl), a);
    if (l <= sl && sl < r) root->x = ActedMonoid::act(root->x, a, 1);
    u32 k = 1 + sl;
    if (k < r) root->r = apply_rec(root->r, max(k, l) - k, r - k, a);
    update(root);
    return root;
  }

  template <typename F>
  pair<np, np> split_max_right_rec(np root, F check, X &x) {
    if (!root) return {nullptr, nullptr};
    prop(root);
    root = copy_node(root);
    X y = Monoid_X::op(x, root->prod);
    if (check(y)) {
      x = y;
      return {root, nullptr};
    }
    np left = root->l, right = root->r;
    if (left) {
      X y = Monoid_X::op(x, root->l->prod);
      if (!check(y)) {
        auto [n1, n2] = split_max_right_rec(left, check, x);
        root->l = n2;
        update(root);
        return {n1, root};
      }
      x = y;
    }
    y = Monoid_X::op(x, root->x);
    if (!check(y)) {
      root->l = nullptr;
      update(root);
      return {left, root};
    }
    x = y;
    auto [n1, n2] = split_max_right_rec(right, check, x);
    root->r = n1;
    update(root);
    return {root, n2};
  }
};