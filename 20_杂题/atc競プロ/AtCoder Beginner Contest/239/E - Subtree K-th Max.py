# 给定一棵树，点有点权，1 号点为根，每次询问以某个点为根的子树内第 k 大的点权是多少。
# n<=1e5 k<=20
# !1.区间第k大 dfs序+主席树/莫队
# !2.因为k很小 所以可以直接暴力维护
from heapq import nlargest
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


def main() -> None:
    def dfs(cur: int, pre: int) -> None:
        dp[cur].append(values[cur])
        for next in adjList[cur]:
            if next == pre:
                continue
            dfs(next, cur)
            dp[cur].extend(dp[next])
        dp[cur] = nlargest(20, dp[cur])

    n, q = map(int, input().split())
    values = list(map(int, input().split()))
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)

    dp = [[] for _ in range(n)]
    dfs(0, -1)
    for _ in range(q):
        root, k = map(int, input().split())
        print(dp[root - 1][k - 1])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
