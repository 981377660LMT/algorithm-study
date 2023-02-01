# 以图判链
# 1. 是一棵树(# 1. 边数 = 节点数 - 1 # 2. 无环)
# 2. 有且仅有两个节点的度为1, 其余节点的度为2


from typing import List
from UnionFind import UnionFindArray


def validChain(n: int, edges: List[List[int]]) -> bool:
    if n != len(edges) + 1:
        return False

    uf = UnionFindArray(n)
    deg = [0] * n
    for u, v in edges:
        if uf.isConnected(u, v):
            return False
        uf.union(u, v)
        deg[u] += 1
        deg[v] += 1

    return deg.count(1) == 2 and deg.count(2) == n - 2


if __name__ == "__main__":
    # C - Path Graph(パスグラフ:链)
    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append([u, v])

    isChain = validChain(n, edges)
    print("Yes" if isChain else "No")
