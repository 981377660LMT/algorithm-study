from Rerooting import Rerooting

from typing import List, Tuple


# https://www.acwing.com/problem/content/289/
def 积蓄程度(n: int, edges: List[Tuple[int, int, int]]) -> List[int]:
    E = int

    def e(root: int) -> E:
        return 0

    def op(childRes1: E, childRes2: E) -> E:
        return childRes1 + childRes2

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        from_ = parent if direction == 1 else cur
        to = cur if direction == 1 else parent
        if fromRes == 0:
            return weight[from_][to]
        return min(fromRes, weight[from_][to])

    R = Rerooting(n)
    weight = [dict() for _ in range(n)]
    for u, v, w in edges:
        R.addEdge(u, v)
        weight[u][v] = w
        weight[v][u] = w

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    return dp


if __name__ == "__main__":
    import sys

    input = sys.stdin.readline

    T = int(input())
    for _ in range(T):
        n = int(input())
        edges = []
        for _ in range(n - 1):
            u, v, w = map(int, input().split())
            u, v = u - 1, v - 1
            edges.append((u, v, w))
        res = 积蓄程度(n, edges)
        print(max(res))
