#pragma once
// select(i,j,k) は (i,j) と (i,k) のうち選ぶ方（j or k）
template <typename F>
vc<int> SMAWK(int H, int W, F select) {
  auto dfs = [&](auto& dfs, vc<int> X, vc<int> Y) -> vc<int> {
    int N = len(X);
    if (N == 0) return {};
    vc<int> YY;
    for (auto&& y: Y) {
      while (len(YY)) {
        int py = YY.back(), x = X[len(YY) - 1];
        if (select(x, py, y) == py) break;
        YY.pop_back();
      }
      if (len(YY) < len(X)) YY.eb(y);
    }
    vc<int> XX;
    FOR(i, 1, len(X), 2) XX.eb(X[i]);
    vc<int> II = dfs(dfs, XX, YY);
    vc<int> I(N);
    FOR(i, len(II)) I[i + i + 1] = II[i];
    int p = 0;
    FOR(i, 0, N, 2) {
      int LIM = (i + 1 == N ? Y.back() : I[i + 1]);
      int best = Y[p];
      while (Y[p] < LIM) {
        ++p;
        best = select(X[i], best, Y[p]);
      }
      I[i] = best;
    }
    return I;
  };
  vc<int> X(H), Y(W);
  iota(all(X), 0), iota(all(Y), 0);
  return dfs(dfs, X, Y);
}