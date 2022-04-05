# Best Currency Path-最佳货币路径
# 求起点到终点的最长路径 如果有环 返回-1 如果到不了 返回0
from collections import defaultdict
from typing import DefaultDict, List, Tuple, TypeVar


T = TypeVar('T', int, str)


def bellman_ford(
    n: int, edges: List[Tuple[T, T, int]], start: T
) -> Tuple[bool, DefaultDict[T, int]]:
    """bellman_ford判断负权环 并求单源最短路"""

    dist = defaultdict(int)  # 起点s到各个点的距离
    dist[start] = 1

    # 松弛i次:其中第i次(i>=1)的内涵为此时至少优化过了过了i-1个`中转点`，最后一次优化了n-1个中转点(即所有点都经过了)
    for _ in range(n):
        for u, v, w in edges:
            if dist[u] * w > dist[v]:
                dist[v] = dist[u] * w

    for u, v, w in edges:
        if dist[u] * w > dist[v]:
            return False, dist  # 存在负权边

    return True, dist


# n,m<=2000
class Solution:
    def solve(
        self, source: str, target: str, sources: List[str], targets: List[str], rates: List[float]
    ):
        n = len(sources)
        edges = []
        for u, v, w in zip(sources, targets, rates):
            edges.append((u, v, w))

        # bellmanFord 检测环 + 求最短路
        isOk, dist = bellman_ford(n, edges, source)
        if not isOk:
            return -1
        res = dist[target]
        return res


print(
    Solution().solve(
        source="CAD",
        target="AUD",
        sources=["CAD", "CAD", "USD", "EUR"],
        targets=["USD", "EUR", "AUD", "AUD"],
        rates=[0.8, 0.5, 2, 3],
    )
)
# We can convert one unit of sources[i] currency to targets[i] currency for rates[i].
