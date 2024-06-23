import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


from collections import defaultdict
from typing import List, Tuple

from typing import Callable, Generic, List, TypeVar

T = TypeVar("T")


class Rerooting(Generic[T]):
    __slots__ = ("adjList", "_n", "_decrement")

    def __init__(self, n: int, decrement: int = 0):
        self.adjList = [[] for _ in range(n)]
        self._n = n
        self._decrement = decrement

    def addEdge(self, u: int, v: int) -> None:
        u -= self._decrement
        v -= self._decrement
        self.adjList[u].append(v)
        self.adjList[v].append(u)

    def rerooting(
        self,
        e: Callable[[int], T],
        op: Callable[[T, T], T],
        composition: Callable[[T, int, int, int], T],
        root=0,
    ) -> List["T"]:
        root -= self._decrement
        assert 0 <= root < self._n
        parents = [-1] * self._n
        order = [root]
        stack = [root]
        while stack:
            cur = stack.pop()
            for next in self.adjList[cur]:
                if next == parents[cur]:
                    continue
                parents[next] = cur
                order.append(next)
                stack.append(next)

        dp1 = [e(i) for i in range(self._n)]
        dp2 = [e(i) for i in range(self._n)]
        for cur in order[::-1]:
            res = e(cur)
            for next in self.adjList[cur]:
                if parents[cur] == next:
                    continue
                dp2[next] = res
                res = op(res, composition(dp1[next], cur, next, 0))
            res = e(cur)
            for next in self.adjList[cur][::-1]:
                if parents[cur] == next:
                    continue
                dp2[next] = op(res, dp2[next])
                res = op(res, composition(dp1[next], cur, next, 0))
            dp1[cur] = res

        for newRoot in order[1:]:
            parent = parents[newRoot]
            dp2[newRoot] = composition(op(dp2[newRoot], dp2[parent]), parent, newRoot, 1)
            dp1[newRoot] = op(dp1[newRoot], dp2[newRoot])
        return dp1


def solve(n: int, edges: List[Tuple[int, int, int]]) -> List[int]:
    E = int

    def e(root: int) -> E:
        return 0

    def op(childRes1: E, childRes2: E) -> E:
        return max(childRes1, childRes2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        return fromRes + weight[cur][parent]

    R = Rerooting(n)
    weight = [defaultdict(int) for _ in range(n)]
    wSum = 0
    for u, v, w in edges:
        R.addEdge(u, v)
        weight[u][v] = w
        weight[v][u] = w
        wSum += w
    dp = R.rerooting(e, op, composition)
    print(dp)
    return [2 * wSum - x for x in dp]


if __name__ == "__main__":
    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v, 1))
    res = solve(n, edges)
    print(*res, sep="\n")
