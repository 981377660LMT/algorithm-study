# Floyd-Warshall 算法求有向图的闭包传递问题。通俗地讲，就是可达性问题。
# O(n^3/w)


class TransitiveClosure:
    """
    Floyd-Warshall 算法求有向图的闭包传递问题，也就是可达性问题。

    2000*2000 => 1s
    3000*3000 => 2.5s
    """

    __slots__ = ("_n", "_canReach", "_hasBuilt")

    def __init__(self, n: int) -> None:
        self._n = n
        self._canReach = [0] * n
        self._hasBuilt = False

    def addDirectedEdge(self, from_: int, to: int) -> None:
        self._hasBuilt = False
        self._canReach[from_] |= 1 << to

    def build(self) -> None:
        if self._hasBuilt:
            return
        self._hasBuilt = True
        n, canReach = self._n, self._canReach
        for k in range(n):
            canReachK = canReach[k]
            for i in range(n):
                if canReach[i] & (1 << k):
                    canReach[i] |= canReachK

    def canReach(self, from_: int, to: int) -> bool:
        if not self._hasBuilt:
            self.build()
        return not not self._canReach[from_] & (1 << to)


if __name__ == "__main__":
    import time
    from typing import List

    n = 3000
    T = TransitiveClosure(n)
    for i in range(n):
        for j in range(n):
            T.addDirectedEdge(i, j)
    time1 = time.time()
    T.build()
    time2 = time.time()
    print(time2 - time1)

    class Solution:
        def checkIfPrerequisite(
            self, numCourses: int, prerequisites: List[List[int]], queries: List[List[int]]
        ) -> List[bool]:
            T = TransitiveClosure(numCourses)
            for from_, to in prerequisites:
                T.addDirectedEdge(from_, to)
            return [T.canReach(from_, to) for from_, to in queries]
