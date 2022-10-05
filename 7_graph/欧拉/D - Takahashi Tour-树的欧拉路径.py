"""
给出一棵有 n 个结点的树，从根节点 1 出发，每次优先遍历序号较小的结点，
如无结点可遍历则返回上一结点，试模拟该过程。

树的欧拉路径/树的欧拉序
"""

import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    def dfs(cur: int, pre: int, path: List[int]) -> None:
        path.append(cur)
        for next in adjList[cur]:
            if next != pre:
                dfs(next, cur, path)
                path.append(cur)

    n = int(input())
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        a, b = map(int, input().split())
        adjList[a - 1].append(b - 1)
        adjList[b - 1].append(a - 1)

    for i in range(n):
        adjList[i].sort()

    res = []
    dfs(0, -1, res)
    print(*[num + 1 for num in res])
