# https://oi-wiki.org/topic/dsu-app/#f
# 给出一棵 n 个点的树，接下来有 m 次操作：
# 加一条从 a_i 到 b_i 的边。
# 询问两个点 u_i 和 v_i 之间是否有至少两条边不相交的路径。

# 每次加边操作，我们就暴力跳并查集。
# 覆盖了一条边后，将这条边对应结点的 f 与父节点合并。
# 这样，每条边至多被覆盖一次，总复杂度 O(n\log n)。
# 使用按秩合并的并查集同样可以做到 O(n\alpha(n))。


from collections import defaultdict
from typing import Callable, DefaultDict, List, Optional, Tuple


class UnionFindRangeOnTree:
    __slots__ = ("part", "_n", "_data", "_treeParents")

    def __init__(self, n: int, treeParents: List[int]):
        self.part = n
        self._n = n
        self._data = [-1] * n
        self._treeParents = treeParents

    def union(
        self, parent: int, child: int, f: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        """将child结点合并到parent结点上,返回是否合并成功."""
        parent, child = self.find(parent), self.find(child)
        if parent == child:
            return False
        self._data[parent] += self._data[child]
        self._data[child] = parent
        self.part -= 1
        if f is not None:
            f(parent, child)
        return True

    def unionRange(
        self, ancestor: int, child: int, f: Optional[Callable[[int, int], None]] = None
    ) -> int:
        """定向合并从祖先ancestor到子孙child路径上的所有节点,返回合并次数."""
        target = self.find(ancestor)
        mergeCount = 0
        while True:
            child = self.find(child)
            if child == target:
                break
            self.union(self._treeParents[child], child, f)
            mergeCount += 1
        return mergeCount

    def find(self, key: int) -> int:
        if self._data[key] < 0:
            return key
        self._data[key] = self.find(self._data[key])
        return self._data[key]

    def isConnected(self, key1: int, key2: int) -> bool:
        return self.find(key1) == self.find(key2)

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for i in range(self._n):
            root = self.find(i)
            groups[root].append(i)
        return groups


def minimumRechableCity(
    n: int,
    parents: List[int],
    queries: List[Tuple[int, int, int]],
) -> List[int]:
    def f(parent: int, child: int) -> None:
        groupMin[parent] = min(groupMin[parent], groupMin[child])

    uf = UnionFindRangeOnTree(n, parents)  # 并查集维护强连通分量
    res = []
    groupMin = list(range(n))
    for op, *args in queries:
        if op == 1:
            child, ancestor = args
            uf.unionRange(ancestor, child, f)
        else:
            x = args[0]
            root = uf.find(x)
            res.append(groupMin[root])
    return res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    sys.setrecursionlimit(int(1e6))

    n = int(input())
    parents = [-1] + [int(x) - 1 for x in input().split()]
    q = int(input())
    queries = []
    for _ in range(q):
        op, *args = map(int, input().split())
        if op == 1:
            queries.append((op, args[0] - 1, args[1] - 1))
        else:
            queries.append((op, args[0] - 1, 0))

    res = minimumRechableCity(n, parents, queries)
    for r in res:
        print(r + 1)
