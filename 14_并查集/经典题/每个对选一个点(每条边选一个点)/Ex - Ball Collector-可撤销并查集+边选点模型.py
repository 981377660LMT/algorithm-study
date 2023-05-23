# https://atcoder.jp/contests/abc302/tasks/abc302_h
# 给定一棵树，每个点有两个值。
# 对于v=1,2,3,...,n，问从点1到点 v的最短路径途径的每个点中，
# 各选一个数，其不同数的个数的最大值。


from SelectOneFromEachPair import SelectOneFromEachPairMap

from typing import List, Tuple


def ballCollector(n: int, edges: List[Tuple[int, int]], pairs: List[Tuple[int, int]]) -> List[int]:
    def dfs(cur: int, pre: int) -> None:
        a, b = pairs[cur]
        uf.union(a, b)
        res[cur] = uf.solve()
        for next in adjList[cur]:
            if next != pre:
                dfs(next, cur)
        uf.undo()

    uf = SelectOneFromEachPairMap()
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
    res = [0] * n
    dfs(0, -1)
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    pairs = [tuple(map(int, input().split())) for _ in range(n)]
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))
    res = ballCollector(n, edges, pairs)
    print(*res[1:])
