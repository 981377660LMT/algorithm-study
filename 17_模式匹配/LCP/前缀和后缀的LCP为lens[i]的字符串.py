# RangeParallelUnionfind-并行合并的并查集

from collections import defaultdict
from typing import DefaultDict, List, Callable


class UnionFindRangeParallel:
    """并行合并的并查集."""

    __slots__ = ("_n", "_queues")

    def __init__(self, n: int) -> None:
        self._n = n
        self._queues = [[] for _ in range(n + 1)]

    def unionParallelly(self, a: int, b: int, len: int) -> None:
        """并行合并[(a,b),(a+1,b+1),...,(a+len-1,b+len-1)]."""
        if len == 0:
            return
        min_ = len if len < self._n else self._n
        self._queues[min_].append((a, b))

    def build(self) -> "UnionFindArray":
        res = UnionFindArray(self._n)
        queue, nextQueue = [], []
        for di in range(self._n, 0, -1):
            queue += self._queues[di]
            nextQueue = []
            for p in queue:
                if res.isConnected(p[0], p[1]):
                    continue
                res.union(p[0], p[1])
                nextQueue.append((p[0] + 1, p[1] + 1))
            queue, nextQueue = nextQueue, queue
        return res


class UnionFindArray:
    __slots__ = ("n", "part", "_parent", "_rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self._parent = list(range(n))
        self._rank = [1] * n

    def find(self, x: int) -> int:
        while x != self._parent[x]:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self._rank[rootX] > self._rank[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        return True

    def unionWithCallback(self, x: int, y: int, f: Callable[[int, int], None]) -> bool:
        """
        f: 合并后的回调函数, 入参为 (big, small)
        """
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self._rank[rootX] > self._rank[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        f(rootY, rootX)
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self._parent))

    def getSize(self, x: int) -> int:
        return self._rank[self.find(x)]

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


if __name__ == "__main__":
    # https://atcoder.jp/contests/yahoo-procon2018-final/tasks/yahoo_procon2018_final_d
    # !前缀和后缀的LCP为lens[i]的字符串
    # !LCP => 并查集
    # 给定长为n的数组lens, 问是否存在一个长度为s的字符串,满足:
    # !s[0:i+1] 和 s[n-(i+1):n] 的最长公共前缀为 lens[i] (0<=i<n)
    # n<=3e5 0<=lens[i]<=i+1

    n = int(input())
    lens = list(map(int, input().split()))
    P = UnionFindRangeParallel(n)
    for i in range(n):
        P.unionParallelly(0, n - (i + 1), lens[i])  # 各个位置的字符相同
    uf = P.build()
    for i in range(n):
        if lens[i] == i + 1:
            continue
        # !s[len[i]]!=s[n-(i+1)+len[i]] (因为前后缀LCP只有len[i])
        if uf.isConnected(lens[i], n - (i + 1) + lens[i]):
            print("No")
            exit(0)
    print("Yes")
