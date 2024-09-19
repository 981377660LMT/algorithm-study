

int LIM = 200100;

struct T {
  int idx;
  ll y1, y2;
};
void solve() {
  LL(N);
  VEC(pi, stone, N);
  for (auto& [x, y]: stone) ++x, ++y;
  vvc<int> XtoY(LIM);
  for (auto& [x, y]: stone) XtoY[x].eb(y);

  // X -> [y1,y2]
  vvc<T> seg(LIM);
  int n = 0;
  FOR(x, LIM) {
    auto Y = XtoY[x];
    sort(all(Y));
    ll a = 0;
    for (auto& y: Y) {
      // [a,y-1]
      if (a < y) seg[x].eb(T{n++, a, y - 1});
      a = y + 1;
    }
    seg[x].eb(T{n++, a, LIM - 1});
  }

  UnionFind uf(n);
  FOR(x, LIM - 1) {
    auto A = seg[x];
    auto B = seg[x + 1];
    while (len(A) && len(B)) {
      auto& [i, y1, y2] = A.back();
      auto& [j, y3, y4] = B.back();
      if (y1 > y4) {
        POP(A);
        continue;
      }
      if (y3 > y2) {
        POP(B);
        continue;
      }
      assert(max(y1, y3) <= min(y2, y4));
      uf.merge(i, j);
      if (y1 < y3) {
        POP(B);
      } else {
        POP(A);
      }
    }
  }
  ll ANS = 0;
  FOR(x, LIM) {
    for (auto& [i, y1, y2]: seg[x]) {
      if (uf[i] != uf[0]) { ANS += y2 - y1 + 1; }
    }
  }
  print(ANS);
}
