"""任意两点距离为直径的集合个数

给定一棵n个节点的树,设它的直径是D,问有多少个集合满足集合中每两个点的距离都为D。(至少染色两个点)
n<=2e5

直径的性质:
!树的每一条直径一定都经过一个公共点 / 一条公共边。
经过的是点还是边取决于直径的长度是奇数还是偶数。

https://www.cnblogs.com/xsl19/p/abc221f.html
https://www.cnblogs.com/Kanoon/p/15380590.html#f---diameter-set
"""

from 树的直径 import calDiameter3


from typing import List, Set
from collections import deque


def bfs(adjList: List[Set[int]], start: int, depth: int) -> int:
    """统计距离start为depth的点的个数"""
    queue = deque([start])
    visited = set([start])
    while queue:
        if depth == 0:
            return len(queue)
        len_ = len(queue)
        for _ in range(len_):
            cur = queue.popleft()
            for next in adjList[cur]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)
        depth -= 1
    return 0


def solve1(adjList: List[Set[int]]) -> int:
    """直径长度为奇数,移去中间边(桥),答案为两边子树内f1*f2"""
    u, v = path[diameter // 2], path[diameter // 2 + 1]
    adjList[u].remove(v)
    adjList[v].remove(u)
    res1 = bfs(adjList, u, diameter // 2)
    res2 = bfs(adjList, v, diameter // 2)
    return res1 * res2 % MOD


def solve2(adjList: List[Set[int]]) -> int:
    """直径长度为偶数,移去中间点(割点),从所有染色情况中减去只染一个或不染的情况
    答案为 (f1+1)*(f2+1)*...*(fn+1)-(f1+f2+...+fn+1)
    """
    u = path[diameter // 2]
    nexts = adjList[u].copy()
    adjList[u].clear()
    for next in nexts:
        adjList[next].remove(u)
    nextRes = [bfs(adjList, next, diameter // 2 - 1) for next in nexts]
    res = 1
    for num in nextRes:
        res = res * (num + 1) % MOD
    res -= sum(nextRes) + 1
    return res % MOD


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353
    INF = int(4e18)
    n = int(input())
    adjList = [set() for _ in range(n)]
    for _ in range(n - 1):
        a, b = map(int, input().split())
        a, b = a - 1, b - 1
        adjList[a].add(b)
        adjList[b].add(a)

    path = calDiameter3(adjList)
    diameter = len(path) - 1
    if diameter % 2 == 0:
        print(solve2(adjList))
    else:
        print(solve1(adjList))
