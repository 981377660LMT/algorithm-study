// // https://uoj.ac/problem/891
// template <typename T, typename F>
// vi nth_element_from_sorted_matrix(ll N, ll M, ll K, F f, int k1 = 0, int k2 = 0,
//                                   bool tr = false) {
//   if (K == 0) return vi(N, 0);
//   if (N > M) {
//     vi A = nth_element_from_sorted_matrix<T>(M, N, K, f, k2, k1, !tr);
//     vi B(N + 1);
//     FOR(i, M) B[0] += 1, B[A[i]] -= 1;
//     FOR(i, N) B[i + 1] += B[i];
//     B.pop_back();
//     return B;
//   }
//   vi A(N);
//   if (K > N) {
//     A = nth_element_from_sorted_matrix<T>(N, M / 2, (K - N) / 2, f, k1, k2 + 1,
//                                           tr);
//     for (auto &a: A) a *= 2;
//     K = K - (K - N) / 2 * 2;
//   }
//   pqg<pair<T, int>> que;
//   auto g = [&](ll i, ll j) -> T {
//     i = ((i + 1) << k1) - 1;
//     j = ((j + 1) << k2) - 1;
//     return (tr ? f(j, i) : f(i, j));
//   };
//   if (A[0] < M) que.emplace(g(0, A[0]), 0);
//   FOR(i, 1, N) if (A[i] < A[i - 1]) que.emplace(g(i, A[i]), i);
//   while (K) {
//     --K;
//     auto [x, i] = POP(que);
//     A[i]++;
//     if (K == 0) break;
//     if (A[i] < M && (i == 0 || A[i - 1] > A[i])) { que.emplace(g(i, A[i]), i); }
//     if (i + 1 < N && A[i + 1] == A[i] - 1) {
//       que.emplace(g(i + 1, A[i + 1]), i + 1);
//     }
//   }
//   return A;
// }

package main

func main() {

}
