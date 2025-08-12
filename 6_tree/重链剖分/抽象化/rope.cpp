#include <bits/stdc++.h>
using namespace std;

template <typename Node>
struct WBLTBase {
  using DATA_T = decltype(Node::dat);
  using SIZE_T = decltype(Node::size);
  vector<Node> n;
  int top;
  static constexpr double ALPHA = 0.292;
  bool too_heavy(SIZE_T x, SIZE_T y) { return y < ALPHA * (x + y); }

  WBLTBase(int size) : n(size), top(1) {}

  const Node& operator[](int x) { return n[x]; }

  template <typename... Args>
  int alloc(Args... args) {
    n[top] = Node(args...);
    return top++;
  }
  bool almost_full() const { return top >= int(n.size()) * 0.99; }

  pair<int, int> cut(int t) {
    n[t].push(n.data());
    return {n[t].l, n[t].r};
  }
  int join(int l, int r) {
    int t = alloc(l, r);
    n[t].pull(n.data());
    return t;
  }

  int merge(int x, int y) {
    if (!x || !y) return x + y;
    if (too_heavy(n[x].size, n[y].size)) {
      auto [a, b] = cut(x);
      if (too_heavy(n[b].size + n[y].size, n[a].size)) {
        auto [c, d] = cut(b);
        return merge(merge(a, c), merge(d, y));
      } else {
        return merge(a, merge(b, y));
      }
    } else if (too_heavy(n[y].size, n[x].size)) {
      auto [a, b] = cut(y);
      if (too_heavy(n[x].size + n[a].size, n[b].size)) {
        auto [c, d] = cut(a);
        return merge(merge(x, c), merge(d, b));
      } else {
        return merge(merge(x, a), b);
      }
    } else {
      return join(x, y);
    }
  }
  pair<int, int> split(int x, SIZE_T k) {
    if (!x or !k) return {0, x};
    if (k >= n[x].size) return {x, 0};
    auto [a, b] = cut(x);
    if (k <= n[a].size) {
      auto [l, r] = split(a, k);
      return {l, merge(r, b)};
    } else {
      auto [l, r] = split(b, k - n[a].size);
      return {merge(a, l), r};
    }
  }
  int shrink(int x, SIZE_T k) {
    if (!x or !k) return 0;
    if (k >= n[x].size) return x;
    auto [a, b] = cut(x);
    if (k <= n[a].size) return shrink(a, k);
    return merge(a, shrink(b, k - n[a].size));
  }
  int k_times(int x, SIZE_T k) {
    if (k == 0) return 0;
    if (k == 1) return x;
    if (k % 2 == 0) {
      int half = k_times(x, k / 2);
      return join(half, half);
    }
    auto dfs = [&](auto rc, SIZE_T s) -> pair<int, int> {
      if (s == 2) return {x, join(x, x)};
      if (s % 2 == 0) {
        auto [a, b] = rc(rc, s / 2);
        return {join(b, a), join(b, b)};
      } else {
        auto [a, b] = rc(rc, (s + 1) / 2);
        return {join(a, a), join(b, a)};
      }
    };
    return dfs(dfs, k).second;
  }

  // 0-indexed
  DATA_T get_kth(int x, SIZE_T k) {
    if (!n[x].l) return n[x].dat;
    if (k < n[n[x].l].size) return get_kth(n[x].l, k);
    return get_kth(n[x].r, k - n[n[x].l].size);
  }
};

namespace RopeImpl {
template <typename DATA_T, typename SIZE_T>
struct Node {
  int l;
  union {
    int r;
    DATA_T dat;
  };
  SIZE_T size;
  Node() = default;
  Node(DATA_T _dat) : l(0), dat(_dat), size(1) {}
  Node(int _l, int _r) : l(_l), r(_r) {}
  void push(Node*) {}
  void pull(Node* n) { size = n[l].size + n[r].size; }
};

template <typename DATA_T, typename SIZE_T = long long>
using Rope = WBLTBase<Node<DATA_T, SIZE_T>>;
}  // namespace RopeImpl
using RopeImpl::Rope;

namespace PersistentSegmentTreeImpl {
template <typename DATA_T, DATA_T (*f)(DATA_T, DATA_T), typename SIZE_T>
struct Node {
  int l, r;
  DATA_T dat;
  SIZE_T size;
  Node() = default;
  Node(DATA_T _dat) : l(0), dat(_dat), size(1) {}
  Node(int _l, int _r) : l(_l), r(_r) {}
  void push(Node*) {}
  void pull(Node* n) {
    dat = f(n[l].dat, n[r].dat);
    size = n[l].size + n[r].size;
  }
};

template <typename DATA_T, DATA_T (*f)(DATA_T, DATA_T), DATA_T (*e)(),
          typename SIZE_T = long long>
struct PersistentSegmentTree : WBLTBase<Node<DATA_T, f, SIZE_T>> {
  using Base = WBLTBase<Node<DATA_T, f, SIZE_T>>;
  using Base::n;

  PersistentSegmentTree(int size) : Base(size) {}

  DATA_T fold(int x, SIZE_T l, SIZE_T r) {
    if (!x or l >= r) return e();
    if (!n[x].l) return (l <= 0 and 0 < r) ? n[x].dat : e();
    SIZE_T ls = n[n[x].l].size, rs = n[n[x].r].size;
    l = clamp<DATA_T>(l, 0, ls + rs), r = clamp<DATA_T>(r, 0, ls + rs);
    if (l >= r) return e();
    if (r - l == ls + rs) return n[x].dat;
    if (r <= ls) return fold(n[x].l, l, r);
    if (ls <= l) return fold(n[x].r, l - ls, r - ls);
    DATA_T L = fold(n[x].l, l, ls);
    DATA_T R = fold(n[x].r, 0, r - ls);
    return f(L, R);
  }
};

};  // namespace PersistentSegmentTreeImpl
using PersistentSegmentTreeImpl::PersistentSegmentTree;

int main() {
  cin.tie(nullptr);
  ios::sync_with_stdio(false);
  int Q;
  cin >> Q;
  Rope<char> t(1.2e7);
  vector<int> n(Q + 2);
  n[0] = t.alloc('0'), n[1] = t.alloc('1');
  long long BIG = 1e18;
  for (int i = 2; i < Q + 2; i++) {
    long long L, R, X;
    cin >> L >> R >> X;
    if (t[n[L]].size >= BIG) {
      n[i] = n[L];
    } else {
      n[i] = t.merge(n[L], n[R]);
      if (t[n[i]].size > BIG) n[i] = t.shrink(n[i], BIG);
    }
    cout << t.get_kth(n[i], X - 1) << "\n";
  }
}
