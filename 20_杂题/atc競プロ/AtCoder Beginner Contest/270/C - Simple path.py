"""dfs寻找路径"""

import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, x, y = map(int, input().split())
    x, y = x - 1, y - 1
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)

    # !dfs要用python3.8交 才不会TLE
    # def dfs(cur: int, pre: int, target: int, path: List[int]) -> None:
    #     if cur == target:
    #         print(*[num + 1 for num in path])
    #         exit(0)
    #     for next in adjList[cur]:
    #         if next == pre:
    #             continue
    #         path.append(next)
    #         dfs(next, cur, target, path)
    #         path.pop()
    #
    # dfs(x, -1, y, [x])

    def stackDfs(adjList: List[List[int]], start: int, target: int) -> List[int]:
        """栈实现dfs,求出start到target的路径"""
        n = len(adjList)
        stack = [(start, -1, 0)]  # cur, pre, dep
        path = [0] * n  # !记录路径(每个深度对应的结点)
        while stack:
            cur, pre, dep = stack.pop()
            path[dep] = cur

            # !处理当前结点的逻辑
            if cur == target:
                return path[: dep + 1]

            for next in adjList[cur]:
                if next == pre:
                    continue
                stack.append((next, cur, dep + 1))

        return []

    path = stackDfs(adjList, x, y)
    print(*[num + 1 for num in path])
