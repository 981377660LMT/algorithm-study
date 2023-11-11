# https://atcoder.jp/contests/abc165/tasks/abc165_f

from bisect import bisect_left
from typing import List

INF = int(1e18)


def LisOnTree(n: int, edges: List[List[int]], values: List[int]) -> List[int]:
    """树上的LIS问题,返回从根节点出发到每个顶点的严格最长上升子序列的长度。"""

    def dfs(cur: int, pre: int) -> None:
        pos = bisect_left(lis, values[cur])
        if pos == len(lis):
            lis.append(values[cur])
            history.append((-1, -1))
        else:
            tmp = lis[pos]
            history.append((pos, tmp))
            lis[pos] = values[cur]

        res[cur] = len(lis)

        for next in adjList[cur]:
            if next == pre:
                continue
            dfs(next, cur)

        pos, tmp = history.pop()
        if pos == -1:
            lis.pop()
        else:
            lis[pos] = tmp

    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)

    lis = []  # lis[i] 表示长度为 i 的上升子序列的最小末尾值
    history = []  # (pos,preLisValue)
    res = [0] * n
    dfs(0, -1)
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    values = list(map(int, input().split()))
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u -= 1
        v -= 1
        edges.append([u, v])
    print(*LisOnTree(n, edges, values), sep="\n")
