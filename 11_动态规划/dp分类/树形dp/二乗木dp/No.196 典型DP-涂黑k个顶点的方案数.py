"""
给定一棵n个点的树,每个顶点染黑或者白
求涂黑k个顶点的方案数模1e9+7,且满足黑色顶点的子树也是黑色
n<=2000 0<=k<=n

!dp[i][j] 表示以i为根的子树中,涂黑j个顶点的方案数
"""

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

if __name__ == "__main__":
    n, k = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        u, v = map(int, input().split())
        adjList[u].append(v)
        adjList[v].append(u)

    def dfs(cur: int, pre: int) -> None:
        subSize[cur] = 1
        dp[cur] = [1, 1]
        for next in adjList[cur]:
            if next == pre:
                continue
            dfs(next, cur)
            merged = [0] * (subSize[cur] + subSize[next] + 1)
            for i in range(subSize[cur] + 1):
                for j in range(subSize[next] + 1):
                    if i != subSize[cur]:  # 不涂黑当前节点
                        merged[i + j] += dp[cur][i] * dp[next][j]
                        merged[i + j] %= MOD
            subSize[cur] += subSize[next]
            dp[cur] = merged
        dp[cur][-1] = 1  # 涂黑当前节点

    subSize = [0] * n
    dp = [[] for _ in range(n)]
    dfs(0, -1)
    print(dp[0][k])
