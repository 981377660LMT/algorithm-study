# https://nyaannyaan.github.io/library/shortest-path/bellman-ford.hpp


from typing import List, Tuple

INF = int(1e18)


def bellmanFord(n: int, edges: List[Tuple[int, int, int]], start: int, target=-1) -> List[int]:
    """
    BellmanFord算法`O(VE)`求解带负权边的单源最短路.

    Args:
        edges (List[Tuple[int, int, int]]): 有向边(from,to,cost).

    Returns:
        Tuple[List[int]]: 起点到各点的最短距离.
        target不存在时, 如果存在负环, 返回空列表.
        target存在时, 如果start到target经过负环, 返回空列表.
    """
    dist = [INF] * n
    dist[start] = 0
    for _ in range(n):  # !中转点数i=0,1,...,n-1
        updated = False
        for from_, to, cost in edges:
            if dist[from_] == INF:
                continue
            cand = dist[from_] + cost
            if cand < dist[to]:
                updated = True
                dist[to] = cand
        if not updated:
            return dist

    # 此时存在负环
    onNegative = [False] * n  # 在负环上的点
    for _ in range(n):
        for from_, to, cost in edges:
            if dist[from_] == INF:
                continue
            cand = dist[from_] + cost
            if cand < dist[to]:
                onNegative[to] = True
                dist[to] = cand
            if onNegative[from_]:
                onNegative[to] = True

    if onNegative[target]:
        return []
    return dist
