void solve() {
  LL(N, M, S, T);
  --S, --T;
  Graph<ll> G(N, 1);
  Graph<ll> Gr(N, 1);
  FOR(_, M) {
    LL(a, b, c);
    --a, --b;
    G.add_edge(a, b, c);
    Gr.add_edge(b, a, c);
  }
  auto [distS, parS] = dijkstra(G, S);
  auto [distT, parT] = dijkstra(Gr, T);
  ll L = distS[T];
 
  // (開始時刻, 終了時刻, 辺番号)
  struct triple {
    ll L, R, id;
  };
  vc<triple> data;
  FORIN(e, G.edges) {
    ll now = distS[e.frm] + distT[e.to] + e.cost;
    if (now == L) data.eb(triple({distS[e.frm], distS[e.frm] + e.cost, e.id}));
  }
  vi X;
  FORIN(t, data) {
    X.eb(t.L);
    X.eb(t.R);
  }
  sort(all(X));
  FORIN(t, data) {
    t.L = LB(X, t.L);
    t.R = LB(X, t.R);
  }
  vi imos(X.size() + 10);
  FORIN(t, data) {
    imos[t.L] += 1;
    imos[t.R] -= 1;
  }
  FOR(i, X.size()) imos[i + 1] += imos[i];
  vc<bool> is_yes(M);
  FORIN(t, data) is_yes[t.id] = (imos[t.L] == 1);
 
  FORIN(e, G.edges) {
    ll now = distS[e.frm] + distT[e.to] + e.cost;
    if (is_yes[e.id]) {
      YES();
    } else {
      ll x = now - L + 1;
      if (x >= e.cost) {
        NO();
      } else {
        print("CAN", x);
      }
    }
  }
}