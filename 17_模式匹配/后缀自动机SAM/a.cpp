void solve() {
  LL(N);
  STR(S);
  Suffix_Array X(S);
  auto [G, dat] = suffix_tree(X);
 
  vc<int> A(N);
  FOR(i, N) A[i] = (S[i] == '(' ? +1 : -1);
  auto Ac = cumsum<int>(A);
  int mi = MIN(Ac);
  for (auto& x: Ac) x -= mi;
  vvc<int> IDS(MAX(Ac) + 1);
  FOR(i, N + 1) IDS[Ac[i]].eb(i);
 
  SegTree<Monoid_Min<int>> seg(Ac);
 
  auto calc = [&](int L, int a, int b) -> int {
    int R = seg.max_right([&](auto e) -> bool { return e >= Ac[L]; }, L);
    // 最大長
    int n = R - L - 1;
    int mi = L + a + 1;
    int ma = min(L + b, L + n);
    int x = Ac[L];
    if (mi > ma) return 0;
    return UB(IDS[x], ma) - LB(IDS[x], mi);
  };
 
  FOR(i, N) calc(i, 0, 0);
 
  ll ANS = 0;
  for (auto& [L, R, lo, hi]: dat) { ANS += calc(X.SA[L], lo, hi); }
  print(ANS);
}
 