# 题意：给定一个学习的拓扑序，输出完成科目n-1所需的总时间
# !建立反图 dfs找到需要学习哪些科目

import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def solve(n: int, adjList: List[List[int]], time: List[int], start: int) -> int:
    def dfs(cur: int) -> None:
        if visited[cur]:
            return
        visited[cur] = True
        for next in adjList[cur]:
            dfs(next)

    visited = [False] * n
    dfs(start)
    return sum(time[i] for i in range(n) if visited[i])


if __name__ == "__main__":

    n = int(input())
    adjList = [[] for _ in range(n)]
    time = [0] * n
    for i in range(n):
        t, _, *nexts = map(int, input().split())
        time[i] = t
        for next in nexts:
            next -= 1
            adjList[i].append(next)

    print(solve(n, adjList, time, n - 1))
