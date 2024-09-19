#include "ds/my_bitset.hpp"

// 行ベクトルを bitset にする
// (2000, 8000) で 300ms 程度（ABC276H）
vc<My_Bitset> solve_linear(int n, int m, vc<My_Bitset> A, My_Bitset b) {
  using BS = My_Bitset;
  assert(len(b) == n);
  int rk = 0;
  FOR(j, m) {
    if (rk == n) break;
    FOR(i, rk + 1, n) if (A[i][j]) {
      swap(A[rk], A[i]);
      if (b[rk] != b[i]) b[rk] = !b[rk], b[i] = !b[i];
      break;
    }
    if (!A[rk][j]) continue;
    FOR(i, n) if (i != rk) {
      if (A[i][j]) { b[i] = b[i] ^ b[rk], A[i] = A[i] ^ A[rk]; }
    }
    ++rk;
  }
  FOR(i, rk, n) if (b[i]) return {};
  vc<BS> res(1, BS(m));

  vc<int> pivot(m, -1);
  int p = 0;
  FOR(i, rk) {
    while (!A[i][p]) ++p;
    res[0][p] = bool(b[i]), pivot[p] = i;
  }
  FOR(j, m) if (pivot[j] == -1) {
    BS x(m);
    x[j] = 1;
    FOR(k, j) if (pivot[k] != -1 && A[pivot[k]][j]) x[k] = 1;
    res.eb(x);
  }
  return res;
}