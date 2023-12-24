# F - Christmas Present 2
# https://atcoder.jp/contests/abc334/tasks/abc334_f

from MonoQueue import MonoQueue

from typing import List, Tuple

INF = int(1e20)


# void solve() {
#   LL(N, K);
#   LL(SX, SY);
#   vi X(N), Y(N);
#   FOR(i, N) { read(X[i], Y[i]); }
#   FOR(i, N) X[i] -= SX, Y[i] -= SY;

#   vc<Re> dp(N + 1, infty<Re>);
#   dp[0] = 0;

#   auto dist = [&](Re x, Re y) -> Re { return sqrt(x * x + y * y); };

#   vc<Re> A(N), B(N - 1);
#   FOR(i, N) A[i] = dist(X[i], Y[i]);
#   FOR(i, N - 1) B[i] = dist(X[i + 1] - X[i], Y[i + 1] - Y[i]);
#   auto Bc = cumsum<Re>(B);

#   SegTree<Monoid_Min<Re>> seg(N + 1);

#   FOR(j, N + 1) {
#     if (j > 0) {
#       int l = j - K;
#       chmax(l, 0);
#       chmin(dp[j], Bc[j - 1] + A[j - 1] + seg.prod(l, j));
#       /*
#       FOR(i, j) {
#         if (j - i > K) continue;
#         chmin(dp[j], dp[i] + Bc[j - 1] - Bc[i] + A[i] + A[j - 1]);
#       }
#       */
#     }
#     if (j < N) seg.set(j, dp[j] + A[j] - Bc[j]);
#   }

#   // print(dp);
#   Re ANS = dp[N];
#   print(ANS);
# }

"""
快递员送货,起始点为(sx, sy),需要到一些房子houses去送货.
送货需要按照顺序送,即先送第一个房子,再送第二个房子,以此类推.
!快递员每次最多携带k个包裹,中途可以回到起点将包裹补充满.
问从起点出发,送完所有房子,回到起点的最短距离是多少.
n,k<=2e5

思路:
1. 每次送有两种选择:
- 从i直接送到i+1
- 从i回到起点,再送到i+1
"""


def christmasPresent2(sx: int, sy: int, houses: List[Tuple[int, int]], k: int) -> float:
    dis


if __name__ == "__main__":
    n, k = map(int, input().split())
    sx, sy = map(int, input().split())
    houses = [tuple(map(int, input().split())) for _ in range(n)]
    print(christmasPresent2(sx, sy, houses, k))  # type: ignore
