# 求一个有向图中有多少个顶点可以到达该图的环
# 反图 + 拓扑排序

from collections import deque
import sys
import os
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    def topoSort(adjList: List[List[int]], indeg: List[int]) -> List[int]:
        """哪些点不会走到环上 最终会抵达稳定点 从稳定点沿着反图拓扑排序"""
        queue = deque([i for i, d in enumerate(indeg) if d == 0])
        while queue:
            cur = queue.popleft()
            for next in adjList[cur]:
                indeg[next] -= 1
                if indeg[next] == 0:
                    queue.append(next)

        return [i for i, d in enumerate(indeg) if d == 0]

    n, m = map(int, input().split())
    rIndeg = [0] * n
    rAdjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        rAdjList[v].append(u)
        rIndeg[u] += 1

    print(n - len(topoSort(rAdjList, rIndeg)))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
