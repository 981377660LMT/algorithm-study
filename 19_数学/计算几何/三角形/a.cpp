
// 点群 A, B を入力 （Point<ll>）
// query(i,j,k)：三角形 AiAjAk 内部の Bl の個数（非負）を返す
// 前計算 O(N^2M)、クエリ O(1)
struct Count_Points_In_Triangles {
  using P = Point<ll>;
  const int LIM = 1'000'000'000 + 10;
  vc<P> A, B;
  vc<int> I, rk; // O から見た偏角ソート順を管理
  vc<int> point; // A[i] と一致する B[j] の数え上げ
  vvc<int> seg;  // 線分 A[i]A[j] 内にある B[k] の数え上げ
  vvc<int> tri;  // OA[i]A[j] 内部にある B[k] の数え上げ
  Count_Points_In_Triangles(vc<P> A, vc<P> B) : A(A), B(B) {
    for (auto&& p: A) assert(-LIM < min(p.x, p.y) && max(p.x, p.y) < LIM);
    for (auto&& p: B) assert(-LIM < min(p.x, p.y) && max(p.x, p.y) < LIM);
    build();
  }

  int query(int i, int j, int k) {
    i = rk[i], j = rk[j], k = rk[k];
    if (i > j) swap(i, j);
    if (j > k) swap(j, k);
    if (i > j) swap(i, j);
    assert(i <= j && j <= k);

    ll d = (A[j] - A[i]).det(A[k] - A[i]);
    if (d == 0) return 0;
    if (d > 0) { return tri[i][j] + tri[j][k] - tri[i][k] - seg[i][k]; }
    int x = tri[i][k] - tri[i][j] - tri[j][k];
    return x - seg[i][j] - seg[j][k] - point[j];
  }

private:
  P take_origin() {
    int N = len(A), M = len(B);
    while (1) {
      P O = P{-LIM, RNG(-LIM, LIM)};
      bool ok = 1;
      FOR(i, N) FOR(j, N) {
        if (A[i] == A[j]) continue;
        if ((A[i] - O).det(A[j] - O) == 0) ok = 0;
      }
      FOR(i, N) FOR(j, M) {
        if (A[i] == B[j]) continue;
        if ((A[i] - O).det(B[j] - O) == 0) ok = 0;
      }
      if (ok) return O;
    }
    return P{};
  }

  void build() {
    P O = take_origin();
    for (auto&& p: A) p = p - O;
    for (auto&& p: B) p = p - O;
    int N = len(A), M = len(B);
    I.resize(N), rk.resize(N);
    iota(all(I), 0);
    sort(all(I), [&](auto& a, auto& b) -> bool { return A[a].det(A[b]) > 0; });
    FOR(i, N) rk[I[i]] = i;
    A = rearrange(A, I);
    point.assign(N, 0);
    seg.assign(N, vc<int>(N));
    tri.assign(N, vc<int>(N));

    FOR(i, N) FOR(j, M) if (A[i] == B[j])++ point[i];
    FOR(i, N) FOR(j, i + 1, N) {
      FOR(k, M) {
        if (A[i].det(B[k]) <= 0) continue;
        if (A[j].det(B[k]) >= 0) continue;
        ll d = (B[k] - A[i]).det(A[j] - A[i]);
        if (d == 0) ++seg[i][j];
        if (d < 0) ++tri[i][j];
      }
    }
  }
};