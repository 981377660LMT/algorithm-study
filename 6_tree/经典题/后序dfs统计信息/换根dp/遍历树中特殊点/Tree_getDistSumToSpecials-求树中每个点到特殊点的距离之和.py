from typing import List, Tuple
from Rerooting import Rerooting

INF = int(1e18)


def getDistSumToSpecials(n: int, edges: List[Tuple[int, int]], specials: List[int]) -> List[int]:
    """求`原树`中每个点到可达的(连通)所有特殊点的距离之和."""
    isSpecial = [False] * n
    for v in specials:
        isSpecial[v] = True

    E = Tuple[int, int]  # (distSum,size) 树中距离之和, 子树中特殊点的大小

    def e(root: int) -> E:
        return (0, 0)

    def op(childRes1: E, childRes2: E) -> E:
        dist1, size1 = childRes1
        dist2, size2 = childRes2
        return (dist1 + dist2, size1 + size2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        preDist, preSize = fromRes
        from_ = cur if direction == 0 else parent
        addedCount = 1 if isSpecial[from_] else 0
        return (preDist + preSize + addedCount, preSize + addedCount)

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)
    dp = R.rerooting(e, op, composition)
    return [d for d, _ in dp]


def getDistSumToSpecialsWeighted(
    n: int, edges: List[Tuple[int, int, int]], specials: List[int]
) -> List[int]:
    isSpecial = [False] * n
    for v in specials:
        isSpecial[v] = True

    E = Tuple[int, int]  # (distSum,size) 树中距离之和, 子树中特殊点的个数

    def e(root: int) -> E:
        return (0, 0)

    def op(childRes1: E, childRes2: E) -> E:
        dist1, size1 = childRes1
        dist2, size2 = childRes2
        return (dist1 + dist2, size1 + size2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        preDist, preSize = fromRes
        from_ = cur if direction == 0 else parent
        addedCount = 1 if isSpecial[from_] else 0
        weight = weights[cur][parent] if cur < parent else weights[parent][cur]
        return (preDist + (preSize + addedCount) * weight, preSize + addedCount)

    R = Rerooting(n)
    weights = [dict() for _ in range(n)]
    for u, v, w in edges:
        if u > v:
            u, v = v, u
        weights[u][v] = w
        R.addEdge(u, v)

    dp = R.rerooting(e, op, composition)
    return [d for d, _ in dp]


if __name__ == "__main__":
    n = 5
    edges = [(0, 1), (1, 2), (1, 3), (3, 4)]
    specials = [0, 3]
    res = getDistSumToSpecials(n, edges, specials)  # type: ignore
    assert res == [2, 2, 4, 2, 4]

    n = 5
    edges = [(0, 1, 2), (1, 2, 3), (1, 3, 4), (3, 4, 5)]
    specials = [0, 3]
    res = getDistSumToSpecialsWeighted(n, edges, specials)  # type: ignore
    assert res == [6, 6, 12, 6, 16]
