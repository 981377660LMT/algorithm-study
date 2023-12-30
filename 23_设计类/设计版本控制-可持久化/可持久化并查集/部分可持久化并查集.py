# 部分可持久化并查集(初始版本为0)
# !只有最新版本可修改,历史版本只能查询

# https://tjkendev.github.io/procon-library/python/union_find/pp_union_find.html
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

from bisect import bisect_right
from collections import defaultdict
from typing import DefaultDict, List


INF = int(1e18)


class PartiallyPersistentUnionFind:
    __slots__ = ("dead", "curVersion", "_history", "_parent", "_size", "_height")

    def __init__(self, n: int):
        self.dead = [INF] * n  # !每个顶点被合并的时刻(不再是每个集合的根的时刻) 0-indexed
        self.curVersion = 0  # !起始版本为0
        self._history = [[(0, 1)] for _ in range(n)]  # (t,s) 表示在t时刻合并完另一个组后的大小s
        self._parent = list(range(n))
        self._size = [1] * n
        self._height = [1] * n

    def union(self, x: int, y: int) -> int:
        """合并x和y所在的集合,返回当前版本号"""
        px = self.find(self.curVersion, x)
        py = self.find(self.curVersion, y)
        if px == py:
            self.curVersion += 1
            return self.curVersion
        if self._height[py] < self._height[px]:
            self._parent[py] = px
            self.dead[py] = self.curVersion
            self._size[px] += self._size[py]
            self._history[px].append((self.curVersion, self._size[px]))  # type: ignore
        else:
            self._parent[px] = py
            self.dead[px] = self.curVersion
            self._size[py] += self._size[px]
            self._history[py].append((self.curVersion, self._size[py]))  # type: ignore
            cand = self._height[px] + 1
            if cand > self._height[py]:
                self._height[py] = cand
        self.curVersion += 1
        return self.curVersion

    def find(self, time: int, u: int) -> int:
        while self.dead[u] < time:
            u = self._parent[u]
        return u

    def isConnected(self, time: int, u: int, v: int) -> bool:
        return self.find(time, u) == self.find(time, v)

    def getSize(self, time: int, u: int) -> int:
        y = self.find(time, u)
        index = bisect_right(self._history[y], (time, INF)) - 1
        return self._history[y][index][1]

    def getGroups(self, time: int) -> "DefaultDict[int, List[int]]":
        groups = defaultdict(list)
        for i in range(len(self._parent)):
            groups[self.find(time, i)].append(i)
        return groups


if __name__ == "__main__":
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
