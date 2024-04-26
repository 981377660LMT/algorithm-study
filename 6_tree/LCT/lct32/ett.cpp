#include <vector>
#include <cassert>
#include <unordered_map>
using namespace std;

namespace titan23 {

  template <class T,
            class F,
            T (*op)(T, T),
            T (*mapping)(F, T),
            F (*composition)(F, F),
            T (*e)(),
            F (*id)()>
  struct EulerTourTree {
    struct Node;
    using NodePtr = Node*;
    int n, group_numbers;
    vector<NodePtr> ptr_vertex;
    unordered_map<long long, NodePtr> ptr_edge;

    struct Node {
      T key, data;
      F lazy;
      NodePtr par, left, right;

      Node() {}

      Node(T key, F lazy) : 
        key(key),
        data(key),
        lazy(lazy),
        par(nullptr),
        left(nullptr),
        right(nullptr) {}
    };

    EulerTourTree(int n) : n(n), group_numbers(n) {
      vector<T> a(n, e());
      _init_build(a);
    }

    EulerTourTree(vector<T> a) : n((int)a.size()), group_numbers((int)a.size()) {
      _init_build(a);
    }

    void _init_build(vector<T> &a) {
      ptr_vertex.resize(n);
      for (int i = 0; i < n; ++i) {
        ptr_vertex[i] = new Node(a[i], id());
      }
    }

    NodePtr _popleft(NodePtr v) {
      v = _left_splay(v);
      if (v->right) v->right->par = nullptr;
      return v->right;
    }

    NodePtr _pop(NodePtr v) {
      v = _right_splay(v);
      if (v->left) v->left->par = nullptr;
      return v->left;
    }

    pair<NodePtr, NodePtr> _split_left(NodePtr v) {
      _splay(v);
      NodePtr x = v, y = v->right;
      if (y) y->par = nullptr;
      x->right = nullptr;
      _update(x);
      return make_pair(x, y);
    }

    pair<NodePtr, NodePtr> _split_right(NodePtr v) {
      _splay(v);
      NodePtr x = v->left, y = v;
      if (x) x->par = nullptr;
      y->left = nullptr;
      _update(y);
      return make_pair(x, y);
    }

    void _merge(NodePtr u, NodePtr v) {
      if ((!u) || (!v)) return;
      u = _right_splay(u);
      _splay(v);
      u->right = v;
      v->par = u;
      _update(u);
    }

    void _splay(NodePtr node) {
      _propagate(node);
      while (node->par && node->par->par) {
        NodePtr pnode = node->par, gnode = pnode->par;
        _propagate(gnode);
        _propagate(pnode);
        _propagate(node);
        node->par = gnode->par;;
        NodePtr tmp1, tmp2;
        if ((gnode->left == pnode) == (pnode->left == node)) {
          if (pnode->left == node) {
            tmp1 = node->right;
            pnode->left = tmp1;
            node->right = pnode;
            pnode->par = node;
            tmp2 = pnode->right;
            gnode->left = tmp2;
            pnode->right = gnode;
            gnode->par = pnode;
          } else {
            tmp1 = node->left;
            pnode->right = tmp1;
            node->left = pnode;
            pnode->par = node;
            tmp2 = pnode->left;
            gnode->right = tmp2;
            pnode->left = gnode;
            gnode->par = pnode;
          }
          if (tmp1) tmp1->par = pnode;
          if (tmp2) tmp2->par = gnode;
        } else {
          if (pnode->left == node) {
            tmp1 = node->right;
            pnode->left = tmp1;
            node->right = pnode;
            tmp2 = node->left;
            gnode->right = tmp2;
            node->left = gnode;
            pnode->par = node;
            gnode->par = node;
          } else {
            tmp1 = node->left;
            pnode->right = tmp1;
            node->left = pnode;
            tmp2 = node->right;
            gnode->left = tmp2;
            node->right = gnode;
            pnode->par = node;
            gnode->par = node;
          }
          if (tmp1) tmp1->par = pnode;
          if (tmp2) tmp2->par = gnode;
        }
        _update(gnode);
        _update(pnode);
        _update(node);
        if (!node->par) return;
        if (node->par->left == gnode) {
          node->par->left = node;
        } else {
          node->par->right = node;
        }
      }
      if (!node->par) return;
      NodePtr pnode = node->par;
      _propagate(pnode);
      _propagate(node);
      if (pnode->left == node) {
        pnode->left = node->right;
        if (pnode->left) pnode->left->par = pnode;
        node->right = pnode;
      } else {
        pnode->right = node->left;
        if (pnode->right) pnode->right->par = pnode;
        node->left = pnode;
      }
      node->par = nullptr;
      pnode->par = node;
      _update(pnode);
      _update(node);
    }

    NodePtr _left_splay(NodePtr node) {
      _splay(node);
      while (node->left) node = node->left;
      _splay(node);
      return node;
    }

    NodePtr _right_splay(NodePtr node) {
      _splay(node);
      while (node->right) node = node->right;
      _splay(node);
      return node;
    }

    void _propagate(NodePtr node) {
      if ((!node) || node->lazy == id()) return;
      if (node->left) {
        node->left->key = mapping(node->lazy, node->left->key);
        node->left->data = mapping(node->lazy, node->left->data);
        node->left->lazy = composition(node->lazy, node->left->lazy);
      }
      if (node->right) {
        node->right->key = mapping(node->lazy, node->right->key);
        node->right->data = mapping(node->lazy, node->right->data);
        node->right->lazy = composition(node->lazy, node->right->lazy);
      }
      node->lazy = id();
    }

    void _update(NodePtr node) {
      _propagate(node->left);
      _propagate(node->right);
      node->data = node->key;
      if (node->left)  node->data = op(node->left->data, node->data);
      if (node->right) node->data = op(node->data, node->right->data);
    }

    // 隣接リストGから構築
    void build(vector<vector<int>> &G) {
      vector<int> seen(n, 0);
      vector<long long> a;
      vector<NodePtr> pool;

      auto dfs = [&] (auto &&dfs, int v, int p) -> void {
        a.emplace_back((long long)v*n+v);
        for (const int &x: G[v]) {
          if (x == p) continue;
          a.emplace_back((long long)v*n+x);
          dfs(dfs, x, v);
          a.emplace_back((long long)x*n+v);
        }
      };

      auto rec = [&] (auto &&rec, int l, int r) -> NodePtr {
        int mid = (l + r) >> 1;
        int u = a[mid]/n, v = a[mid]%n;
        NodePtr node;
        if (u == v) {
          node = ptr_vertex[u];
          seen[u] = 1;
        } else {
          node = new Node(e(), id());
          ptr_edge[a[mid]] = node;
        }

        if (l != mid) {
          node->left = rec(rec, l, mid);
          node->left->par = node;
        }
        if (mid+1 != r) {
          node->right = rec(rec, mid+1, r);
          node->right->par = node;
        }
        _update(node);
        return node;
      };

      for (int root = 0; root < n; ++root) {
        if (seen[root]) continue;
        a.clear();
        dfs(dfs, root, -1);
        rec(rec, 0, (int)a.size());
      }
    }

    // 辺{u, v}を追加する
    void link(const int u, const int v) {
      reroot(u);
      reroot(v);
      NodePtr uv_node = new Node(e(), id());
      NodePtr vu_node = new Node(e(), id());
      ptr_edge[(long long)u*n+v] = uv_node;
      ptr_edge[(long long)v*n+u] = vu_node;
      NodePtr u_node = ptr_vertex[u];
      NodePtr v_node = ptr_vertex[v];
      _merge(u_node, uv_node);
      _merge(uv_node, v_node);
      _merge(v_node, vu_node);
      --group_numbers;
    }

    // 辺{u, v}を削除する
    void cut(const int u, const int v) {
      reroot(v);
      reroot(u);
      NodePtr uv_node = ptr_edge[(long long)u*n+v];
      NodePtr vu_node = ptr_edge[(long long)v*n+u];
      ptr_edge.erase((long long)u*n+v);
      ptr_edge.erase((long long)v*n+u);
      NodePtr a, c, _;
      tie(a, _) = _split_left(uv_node);
      tie(_, c) = _split_right(vu_node);
      a = _pop(a);
      c = _popleft(c);
      _merge(a, c);
      ++group_numbers;
    }

    // 辺{u, v}がなければ追加する
    bool merge(const int u, const int v) {
      if (same(u, v)) return false;
      link(u, v);
      return true;
    }

    // 辺{u, v}があれば削除する
    bool split(const int u, const int v) {
      if (ptr_edge.find((long long)u*n+v) == ptr_edge.end() || ptr_edge.find((long long)v*n+n) == ptr_edge.end()) return false;
      cut(u, v);
      return true;
    }

    // 代表元？
    NodePtr leader(const int v) {
      return _left_splay(ptr_vertex[v]);
    }

    // 根をvにする
    void reroot(const int v) {
      NodePtr node = ptr_vertex[v];
      auto[x, y] = _split_right(node);
      _merge(y, x);
      _splay(node);
    }

    // 連結判定
    bool same(const int u, const int v) {
      NodePtr u_node = ptr_vertex[u];
      NodePtr v_node = ptr_vertex[v];
      _splay(u_node);
      _splay(v_node);
      return (u_node->par != nullptr || u_node == v_node);
    }

    // vを根とする部分木にfを作用、ただしvの親はp(or -1)
    void subtree_apply(const int v, const int p, const F f) {
      NodePtr v_node = ptr_vertex[v];
      reroot(v);
      if (p == -1) {
        _splay(v_node);
        v_node->key = mapping(f, v_node->key);
        v_node->data = mapping(f, v_node->data);
        v_node->lazy = composition(f, v_node->lazy);
        return;
      }
      reroot(p);
      NodePtr a, b, d;
      tie(a, b) = _split_right(ptr_edge[(long long)p*n+v]);
      tie(b, d) = _split_left(ptr_edge[(long long)v*n+p]);
      _splay(v_node);
      v_node->key = mapping(f, v_node->key);
      v_node->data = mapping(f, v_node->data);
      v_node->lazy = composition(f, v_node->lazy);
      _propagate(v_node);
      _merge(a, b);
      _merge(b, d);
    }

    // vを根とする部分木の総和、ただしvの親はp(or -1)
    T subtree_sum(const int v, const int p) {
      NodePtr v_node = ptr_vertex[v];
      reroot(v);
      if (p == -1) {
        _splay(v_node);
        return v_node->data;
      }
      reroot(p);
      NodePtr a, b, d;
      tie(a, b) = _split_right(ptr_edge[(long long)p*n+v]);
      tie(b, d) = _split_left(ptr_edge[(long long)v*n+p]);
      _splay(v_node);
      T res = v_node->data;
      _merge(a, b);
      _merge(b, d);
      return res;
    }

    // 連結成分の個数を返す
    int group_count() const {
      return group_numbers;
    }

    // vの値を取得
    T get_vertex(const int v) {
      NodePtr node = ptr_vertex[v];
      _splay(node);
      return node->key;
    }

    // vの値をvalに変更
    void set_vertex(const int v, const T val) {
      NodePtr node = ptr_vertex[v];
      _splay(node);
      node->key = val;
      _update(node);
    }
  };
}  // namespace titan23
