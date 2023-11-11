# 以图判树
# 1. 一个连通分量
# 2. 连通分量的 边数 = 节点数 - 1 (否则这个连通分量有环)

from collections import defaultdict
from typing import List
from UnionFind import UnionFindArray, UnionFindMap


def isTree(n: int, edges: List[List[int]]) -> bool:
    if n != len(edges) + 1:
        return False
    uf = UnionFindArray(n)
    for u, v in edges:
        uf.union(u, v)
    return uf.part == 1


def isTree2(adjList: List[List[int]], group: List[int]) -> bool:
    """给定无向图的一个连通分量,判断是否是树"""
    edges = sum(len(adjList[u]) for u in group) // 2
    return edges == len(group) - 1


if __name__ == "__main__":
    # https://atcoder.jp/contests/arc111/tasks/arc111_b
    # B - Reversible Cards
    # !给定一些卡片,正反面标有数字,你需要反转一些卡片,使得所有卡片的正面数字种类最多

    # 正反面连边,问题变为选取每个边的一个端点，问最多选取多少个
    # 对每个大小为n的连通块
    # 树: n-1
    # 有环: n
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    adjMap = defaultdict(list)
    uf = UnionFindMap()
    for _ in range(n):
        u, v = map(int, input().split())
        adjMap[u].append(v)
        adjMap[v].append(u)
        uf.union(u, v)

    res = 0
    groups = uf.getGroups()
    for group in groups.values():
        edges = sum(len(adjMap[u]) for u in group) // 2
        noCycle = edges == len(group) - 1  # 连通块无环
        res += len(group) - 1 if noCycle else len(group)
    print(res)
