# https://www.luogu.com.cn/problem/CF1624G
# 给定一张无向图，定义其生成树的边权为 OR 和，求最小生成树的边权 OR 和(或最小生成树)
# 最大边权<=1e9
# !从高到低按位贪心，用并查集判断该位是不是必须取(去掉这位后能否满足存在生成树)，不是的话就不取

from typing import List, Tuple


class UnionFindArraySimple:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(self, key1: int, key2: int) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        if self._data[key] < 0:
            return key
        self._data[key] = self.find(self._data[key])
        return self._data[key]

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]


def cf1624G(n: int, edges: List[Tuple[int, int, int]]) -> int:
    maxWeight = max(w[2] for w in edges)
    maxBit = maxWeight.bit_length()

    res = 0
    for i in range(maxBit, -1, -1):
        uf = UnionFindArraySimple(n)
        for u, v, w in edges:
            if not (w >> i) & 1:
                uf.union(u, v)
        if uf.part > 1:  # 不能形成生成树
            res |= 1 << i
        else:
            edges = [e for e in edges if not (e[2] >> i) & 1]
    return res


if __name__ == "__main__":
    import sys

    input = sys.stdin.readline
    T = int(input())
    for _ in range(T):
        input()
        n, m = map(int, input().split())
        edges = []
        for _ in range(m):
            u, v, w = map(int, input().split())
            u, v = u - 1, v - 1
            edges.append((u, v, w))
        print(cf1624G(n, edges))
