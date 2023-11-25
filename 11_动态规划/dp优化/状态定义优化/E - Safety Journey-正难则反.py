"""
有N座城市,每对城市之间都有一条双向通行的道路,
因为某些原因,现在有M条道路无法使用。
现在从城市1出发旅行K天,最后一天回到城市1。
问一共有多少种不同的旅行方案。

n<=5000 k<=5000 M<=5000
!普通的dp是O(n^2k)的 要想优化到O(5000*5000)

!正难则反
dp[i-1][pre] => dp[i][cur] 转移时
!不计算 `合理的转移次数` 而是 `总次数 - 不合理的转移次数`(往不通的路上靠)

E - Safety Journey
https://atcoder.jp/contests/abc212/submissions/35392362
"""

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m, k = map(int, input().split())
    ban = [[] for _ in range(n)]  # !表示路不通的邻接表
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        ban[u].append(v)
        ban[v].append(u)

    dp = [0] * (n + 10)
    dp[0] = 1
    for i in range(k):
        ndp, sum_ = [0] * (n + 10), sum(dp) % MOD
        for cur in range(n):
            ndp[cur] = (sum_ - dp[cur]) % MOD  # !不能原地转移
            for pre in ban[cur]:
                ndp[cur] = (ndp[cur] - dp[pre]) % MOD  # !不能从不通的路转移过来
        dp = ndp
    print(dp[0])
