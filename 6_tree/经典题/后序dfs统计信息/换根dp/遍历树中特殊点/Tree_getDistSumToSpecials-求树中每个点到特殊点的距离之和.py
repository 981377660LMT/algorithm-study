from typing import List, Tuple
from Rerooting import Rerooting

INF = int(1e18)


def getDistSumToSpecials(n: int, edges: List[Tuple[int, int]], specials: List[int]) -> List[int]:
    """求`原树`中每个点到可达的(连通)所有特殊点的距离之和."""
    weights = [0] * n  # 在这个点处, 有多少个特殊点
    for v in specials:
        weights[v] += 1

    E = Tuple[int, int]  # (distSum,size) 树中距离之和, 子树中特殊点的大小

    def e(root: int) -> E:
        return (0, 0)

    def op(childRes1: E, childRes2: E) -> E:
        dist1, size1 = childRes1
        dist2, size2 = childRes2
        return (dist1 + dist2, size1 + size2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        dist, size = fromRes
        from_ = cur if direction == 0 else parent
        count = weights[from_]
        return (dist + size + count, size + count)

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)
    dp = R.rerooting(e, op, composition)
    return [d for d, _ in dp]


if __name__ == "__main__":
    n = 5
    edges = [(0, 1), (1, 2), (1, 3), (3, 4)]
    specials = [0, 3]
    res = getDistSumToSpecials(n, edges, specials)
    assert res == [2, 2, 4, 2, 4]
