#include "convex/monotone_minima.hpp"

template <typename T>
vc<T> minplus_convolution_convex_convex(vc<T>& A, vc<T>& B) {
  int n = len(A), m = len(B);
  if (n == 0 && m == 0) return {};
  vc<T> C(n + m - 1, infty<T>);
  while (n > 0 && A[n - 1] == infty<T>) --n;
  while (m > 0 && B[m - 1] == infty<T>) --m;
  if (n == 0 && m == 0) return C;
  int a = 0, b = 0;
  while (a < n && A[a] == infty<T>) ++a;
  while (b < m && B[b] == infty<T>) ++b;
  C[a + b] = A[a] + B[b];
  for (int i = a + b + 1; i < n + m - 1; ++i) {
    if (b == m - 1 || (a != n - 1 && A[a + 1] + B[b] < A[a] + B[b + 1])) {
      chmin(C[i], A[++a] + B[b]);
    } else {
      chmin(C[i], A[a] + B[++b]);
    }
  }
  return C;
}

template <typename T>
vc<T> minplus_convolution_arbitrary_convex(vc<T>& A, vc<T>& B) {
  int n = len(A), m = len(B);
  if (n == 0 && m == 0) return {};
  vc<T> C(n + m - 1, infty<T>);
  while (m > 0 && B[m - 1] == infty<T>) --m;
  if (m == 0) return C;
  int b = 0;
  while (b < m && B[b] == infty<T>) ++b;

  auto select = [&](int i, int j, int k) -> bool {
    if (i < k) return false;
    if (i - j >= m - b) return true;
    return A[j] + B[b + i - j] >= A[k] + B[b + i - k];
  };
  vc<int> J = monotone_minima(n + m - b - 1, n, select);
  FOR(i, n + m - b - 1) {
    T x = A[J[i]], y = B[b + i - J[i]];
    if (x < infty<T> && y < infty<T>) C[b + i] = x + y;
  }
  return C;
}

template <typename T, bool convA, bool convB>
vc<T> minplus_convolution(vc<T>& A, vc<T>& B) {
  static_assert(convA || convB);
  if constexpr (convA && convB) return minplus_convolution_convex_convex(A, B);
  if constexpr (convA && !convB)
    return minplus_convolution_arbitrary_convex(B, A);
  if constexpr (convB && !convA)
    return minplus_convolution_arbitrary_convex(A, B);
  return {};
}