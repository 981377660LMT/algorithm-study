# 以图判树
# 1. 边数 = 节点数 - 1
# 2. 无环

from typing import List
from UnionFind import UnionFindArray


def validTree(n: int, edges: List[List[int]]) -> bool:
    if n != len(edges) + 1:
        return False

    uf = UnionFindArray(n)
    for u, v in edges:
        if uf.isConnected(u, v):
            return False
        uf.union(u, v)

    return True
