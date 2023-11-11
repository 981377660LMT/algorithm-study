# 删除顶点后的最大连通分量数
# 给定一棵树,删除其中的任意个顶点,求删除后的最大连通分量数
# https://blog.hamayanhamayan.com/entry/2019/01/03/014943
# !dp[i][0/1]表示以i为根的子树中,不删除/删除i后的最大连通分量数

import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


if __name__ == "__main__":
    n = int(input())
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)

    def dfs(cur: int, pre: int) -> List[int]:
        res = [1, 0]  # 不删除/删除当前节点后的最大连通分量数
        for next in adjList[cur]:
            if next == pre:
                continue
            nextRes = dfs(next, cur)
            res[0] += max(nextRes[0] - 1, nextRes[1])
            res[1] += max(nextRes[0], nextRes[1])
        return res

    res = dfs(0, -1)
    print(max(res))
