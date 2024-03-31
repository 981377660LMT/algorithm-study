/**
 * helper for edge CD: given tree, update node's value, find sum of neighbors' values
 */
template <class T> struct sum_adj {
  int n;
  vector<T> sum, sum_ch;
  vector<int> p;
  /**
   * @param adj undirected, unrooted tree
   * @param a_sum a_sum[u] = initial value for node u
   * @time O(n)
   * @space various O(n) vectors are allocated; recursion stack for dfs is O(n)
   */
  sum_adj(const vector<vector<int>>& adj, const vector<T>& a_sum) : n(ssize(a_sum)), sum(a_sum), sum_ch(n), p(n, -1) {
    auto dfs = [&](auto&& self, int u) -> void {
      for (int v : adj[u])
        if (v != p[u])
          p[v] = u, sum_ch[u] += sum[v], self(self, v);
    };
    dfs(dfs, 0);
  }
  /**
   * @param u node
   * @param delta value to add
   * @time O(1)
   * @space O(1)
   */
  void update(int u, T delta) {
    sum[u] += delta;
    if (p[u] != -1) sum_ch[p[u]] += delta;
  }
  /**
   * @param u node
   * @returns sum of u's neighbors values
   * @time O(1)
   * @space O(1)
   */
  T query(int u) { return sum_ch[u] + (p[u] != -1 ? sum[p[u]] : 0); }
};
