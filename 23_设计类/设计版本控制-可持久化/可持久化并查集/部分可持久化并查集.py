# 部分可持久化并查集(初始版本为0)
# !只有最新版本可修改,历史版本只能查询

# 部分永続Union-Findでは以下の操作を行える。

# union(u,v)
#  : 頂点uの属するグループと要素vの属するグループを1つにまとめる
# find(t,u)
#  : 時刻tにおける頂点uの親頂点を求める
# size(t,u)
#  : 時刻tにおける頂点uの属するグループの頂点数を求める
# 使用:
# uf = PartiallyPersistentUnionFind(10)
# v1 = uf.union(0, 1)  # 合并0和1,返回新版本号1
# print(uf.getGroups(v1))
# print(uf.getSize(v1, 1))  # v1版本的1所在集合的大小为2
# print(uf.isConnected(v1, 0, 1))  # v1版本的0和1是连通的
# print(uf.isConnected(v1, 0, 2))  # v1版本的0和2是不连通的

from collections import defaultdict
from typing import DefaultDict, List


class PartiallyPersistentUnionFind:
    """部分可持久化并查集(初始版本为0)."""

    __slots__ = ("curVersion", "_history", "_data", "_last")

    def __init__(self, n: int):
        self.curVersion = 0
        self._history = [[(0, -1)] for _ in range(n)]
        self._data = [-1] * n
        self._last = [int(1e9)] * n

    def union(self, x: int, y: int) -> int:
        """合并x和y所在的集合,返回当前版本号."""
        self.curVersion += 1
        x, y = self.find(self.curVersion, x), self.find(self.curVersion, y)
        if x == y:
            return self.curVersion
        if self._data[x] > self._data[y]:
            x, y = y, x
        self._data[x] += self._data[y]
        self._history[x].append((self.curVersion, self._data[x]))  # type: ignore
        self._data[y] = x
        self._last[y] = self.curVersion
        return self.curVersion

    def find(self, time: int, u: int) -> int:
        if time < self._last[u]:
            return u
        return self.find(time, self._data[u])

    def isConnected(self, time: int, u: int, v: int) -> bool:
        return self.find(time, u) == self.find(time, v)

    def getSize(self, time: int, u: int) -> int:
        u = self.find(time, u)
        tmp = self._history[u]
        left, right = 0, len(tmp) - 1
        while left <= right:
            mid = (left + right) // 2
            if tmp[mid][0] <= time:
                left = mid + 1
            else:
                right = mid - 1
        return -tmp[left - 1][1]

    def getGroups(self, time: int) -> "DefaultDict[int, List[int]]":
        groups = defaultdict(list)
        for i in range(len(self._data)):
            groups[self.find(time, i)].append(i)
        return groups


if __name__ == "__main__":

    def demo() -> None:
        uf = PartiallyPersistentUnionFind(10)
        print(uf.getSize(0, 1))
        uf.union(0, 1)
        print(uf.getSize(0, 1))
        print(uf.getSize(1, 1))
        uf.union(1, 2)
        print(uf.getSize(1, 1))
        print(uf.getSize(2, 1))
        print(uf.getSize(2, 2))

    def unionSets() -> None:
        # https://atcoder.jp/contests/code-thanks-festival-2017-open/tasks/code_thanks_festival_2017_h
        # 给定n个集合,初始时第i个集合只有一个元素i (i=1,2,...,n)
        # 之后进行m次合并操作,每次合并ai和bi所在的集合
        # 如果ai和bi在同一个集合,则无事发生
        # 给定q个询问,问ai和bi是在第几次操作后第一次连通的,如果不连通则输出-1
        n, m = map(int, input().split())
        uf = PartiallyPersistentUnionFind(n)
        for _ in range(m):
            a, b = map(int, input().split())
            a, b = a - 1, b - 1
            uf.union(a, b)

        q = int(input())
        for _ in range(q):
            a, b = map(int, input().split())
            a, b = a - 1, b - 1
            if not uf.isConnected(uf.curVersion, a, b):
                print(-1)
                continue

            left, right = 0, m  # !二分版本号
            while left <= right:
                mid = (left + right) // 2
                if uf.isConnected(mid, a, b):  # !mid版本是否连通
                    right = mid - 1
                else:
                    left = mid + 1
            print(left)

    def stampRally() -> None:
        #  https://atcoder.jp/contests/agc002/tasks/agc002_d
        #  一张连通图，q 次询问从两个点 x 和 y 出发，
        #  希望经过的点数量等于 z（每个点可以重复经过，但是重复经过只计算一次）
        #  求经过的边最大编号最小是多少。

        import sys

        input = sys.stdin.readline

        n, m = map(int, input().split())
        edges = []
        for _ in range(m):
            a, b = map(int, input().split())
            edges.append((a - 1, b - 1))
        q = int(input())
        queries = []
        for _ in range(q):
            x, y, z = map(int, input().split())
            queries.append((x - 1, y - 1, z))

        uf = PartiallyPersistentUnionFind(n)
        for u, v in edges:
            uf.union(u, v)

        res = [0] * q
        for i, (x, y, z) in enumerate(queries):

            def check(mid: int) -> bool:
                if uf.isConnected(mid, x, y):
                    size = uf.getSize(mid, x)
                    return size >= z
                else:
                    size1, size2 = uf.getSize(mid, x), uf.getSize(mid, y)
                    return size1 + size2 >= z

            left, right = 1, m
            while left <= right:
                mid = (left + right) // 2
                if check(mid):
                    right = mid - 1
                else:
                    left = mid + 1
            res[i] = left

        print(*res, sep="\n")

    stampRally()
