# 起点任意，求访问树中所有点的最长距离


from typing import List, Tuple
from 树的直径 import calDiameter


def maxDistTranverse(n: int, tree: List[List[Tuple[int, int]]], backToStart=False) -> int:
    """tree是无向图邻接表."""
    weightSum = 0
    for cur in range(n):
        for _, weight in tree[cur]:
            weightSum += weight
    weightSum //= 2

    if backToStart:
        return weightSum * 2
    diameter, _ = calDiameter(tree)
    return weightSum * 2 - diameter


if __name__ == "__main__":
    n = 6
    edges = [[1, 2], [1, 3], [3, 4], [3, 5], [5, 6]]
    tree = [[] for _ in range(n)]
    for u, v in edges:
        tree[u - 1].append((v - 1, 1))
        tree[v - 1].append((u - 1, 1))
    assert maxDistTranverse(n, tree) == 6
