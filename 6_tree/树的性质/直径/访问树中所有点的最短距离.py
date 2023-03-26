# 起点任意，求访问树中所有点的最短距离-直径求法(其实更好的解法是换根dp)
# !6_tree/经典题/后序dfs统计信息/换根dp/遍历树中特殊点/访问树上所有顶点的最短路径.py

# !回到起点:边数*2 (边权之和*2)
# !不回到起点:边数*2-直径 (边权之和*2-直径)


from typing import List, Tuple
from 树的直径 import calDiameter


def minDistTranverse(n: int, tree: List[List[Tuple[int, int]]], backToStart=False) -> int:
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
    assert minDistTranverse(n, tree) == 6
