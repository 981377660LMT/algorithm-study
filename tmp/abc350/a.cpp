
void solve() {
  LL(N, Q);
  vc<BS> bs(N);
  vc<set<int>> nbd(N);
  int t = 100;

  auto add = [&](int u, int v) -> void {
    nbd[u].emplace(v);
    if (len(nbd[u]) == t) {
      bs[u].resize(N);
      for (auto& x: nbd[u]) bs[u][x] = 1;
    }
    if (len(bs[u]) == N) bs[u][v] = 1;
  };

  auto find = [&](int u, int v) -> int {
    if (len(nbd[u]) > len(nbd[v])) swap(u, v);
    if (len(nbd[u]) >= t) {
      int x = (bs[u] & bs[v]).next(0);
      return (x == N ? -1 : x);
    }
    if (len(nbd[v]) >= t) {
      for (auto& x: nbd[u])
        if (bs[v][x]) return x;
      return -1;
    }
    for (auto& x: nbd[u])
      if (nbd[v].count(x)) return x;
    return -1;
  };

  ll X = 0;
  FOR(Q) {
    LL(a, b, c);
    ll A = 1 + ((a * (1 + X)) % 998244353) % 2;
    ll B = 1 + ((b * (1 + X)) % 998244353) % N;
    ll C = 1 + ((c * (1 + X)) % 998244353) % N;
    ll u = --B, v = --C;
    if (A == 1) { add(u, v), add(v, u); }
    if (A == 2) {
      int ans = find(u, v) + 1;
      print(ans);
      X = ans;
    }
  }
}
