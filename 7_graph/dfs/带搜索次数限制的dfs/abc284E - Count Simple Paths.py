# https://atcoder.jp/contests/abc284/tasks/abc284_e
# 给定一张图，点度数不超过10。
# 问从1号点（编号从 1开始）的简单路径（不经过重复点）条数 k，
# !输出  min(k,1e6)

from typing import List

MAX = int(1e6)


def countSimplePath(n: int, graph: List[List[int]]) -> int:
    def dfs(cur: int) -> None:
        nonlocal count
        if count >= MAX:  # !剪枝
            return
        if visited[cur]:
            return
        visited[cur] = True
        count += 1
        for next in graph[cur]:
            dfs(next)
        visited[cur] = False

    count = 0
    visited = [False] * n
    dfs(0)
    return count


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, m = map(int, input().split())
    graph = [[] for _ in range(n)]
    for _ in range(m):
        a, b = map(int, input().split())
        graph[a - 1].append(b - 1)
        graph[b - 1].append(a - 1)
    print(countSimplePath(n, graph))
