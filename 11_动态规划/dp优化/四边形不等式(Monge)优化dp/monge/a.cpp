#include "convex/smawk.hpp"

template <typename T, bool concaveA, bool concaveB>
vc<T> maxplus_convolution_concave(vc<T> A, vc<T> B) {
  static_assert(concaveA || concaveB);
  assert(infty<int> < infty<int> + infty<int>);
  if (!concaveB) swap(A, B);
  int NA = len(A), NB = len(B);
  int N = NA + NB - 1;
  int L = 0, R = NB;
  while (L < R && B[L] == -infty<int>) ++L;
  if (L == R) return vc<T>(N, -infty<int>);
  while (B[R - 1] == -infty<int>) --R;
  B = {B.begin() + L, B.begin() + R};
  int nB = R - L;
  int n = NA + nB - 1;

  auto select = [&](int i, int j, int k) -> int {
    if (i < k) return j;
    if (i - j >= nB) return k;
    return (A[j] + B[i - j] < A[k] + B[i - k] ? k : j);
  };

  vc<int> J = SMAWK(n, NA, select);
  vc<T> C(N, -infty<int>);
  FOR(i, n) C[L + i] = A[J[i]] + (A[J[i]] == -infty<int> ? 0 : B[i - J[i]]);
  return C;
}