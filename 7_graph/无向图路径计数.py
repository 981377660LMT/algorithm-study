# E - Count Simple Paths

# 给定一个每个点度数不超过10的无向图,求从0出发的简单路径数
# 如果路径数超过10**6,输出10**6

# 路径计数->dfs求路径
# !dfs最多调用1e6次,每次转移不超过10 次,所以复杂度是O(1e6*10)=O(1e7)


from typing import List


def countSimplePaths(n: int, adjList: List[List[int]], start: int, limit=int(1e6)) -> int:
    def dfs(cur: int) -> None:
        nonlocal res
        if res == limit:
            return

        if visited[cur]:
            return
        visited[cur] = True
        path.append(cur)
        res += 1
        for next in adjList[cur]:
            dfs(next)
        visited[cur] = False
        path.pop()

    path = []
    visited = [False] * n
    res = 0
    dfs(start)
    return res


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


if __name__ == "__main__":

    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)

    print(countSimplePaths(n, adjList, 0))
