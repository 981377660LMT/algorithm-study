from heapq import heappop, heappush
from typing import Tuple, List

INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def minDistToOther(
    adjList: List[List[Tuple[int, int]]], criticals: List[int]
) -> Tuple[List[int], List[int]]:
    """
    给定一个无负权的带权图和一个点集,对点集内的每个点V,求V到点集内其他点距离的最小值,以及到V最近的点是谁.
    !换一种描述方法：给定一张图，有一些黑点和白点，对每个黑点，求出它到其他黑点的最近距离.
    按照 criticals 中点的顺序返回答案.
    """
    n = len(adjList)
    dist = [INF] * n
    source1, source2 = [-1] * n, [-1] * n
    pq = [(0, v, v) for v in criticals]
    while pq:
        curDist, cur, curSource = heappop(pq)
        if curSource == source1[cur] or curSource == source2[cur]:
            continue
        if source1[cur] == -1:
            source1[cur] = curSource
        elif source2[cur] == -1:
            source2[cur] = curSource
        else:
            continue

        if curSource != cur:  # 出发点不为自己时，更新距离
            dist[cur] = min2(dist[cur], curDist)
        for next, weight in adjList[cur]:
            heappush(pq, (curDist + weight, next, curSource))  # type: ignore

    for i, v in enumerate(criticals):
        dist[i] = dist[v]
        source2[i] = source2[v]
    return dist[: len(criticals)], source2[: len(criticals)]


if __name__ == "__main__":
    n = 10
    edges = [
        [0, 1, 4],
        [1, 2, 8],
        [2, 3, 3],
        [3, 4, 3],
        [4, 5, 1],
        [4, 6, 1],
        [3, 7, 4],
        [2, 8, 1],
        [1, 9, 1],
    ]
    points = [5, 6, 7, 8, 9]
    adjList = [[] for _ in range(n)]
    for u, v, w in edges:
        adjList[u].append((v, w))
        adjList[v].append((u, w))
    print(minDistToOther(adjList, points))
