# この問題は問題 G と似た設定です。問題文の相違点を赤字で示します。

# H 行
# W 列のグリッドがあり、グリッドの各マスは赤色あるいは緑色に塗られています。

# グリッドの上から
# i 行目、左から
# j 列目のマスをマス
# (i,j) と表記します。

# マス
# (i,j) の色は文字
# S
# i,j
# ​
#   で表され、
# S
# i,j
# ​
#  = . のときマス
# (i,j) は赤色、
# S
# i,j
# ​
#  = # のときマス
# (i,j) は緑色に塗られています。

# グリッドにおいて、緑色に塗られたマスを頂点集合、隣り合った
# 2 つの緑色のマスを結ぶ辺全体を辺集合としたグラフにおける連結成分の個数を 緑の連結成分数 と呼びます。ただし、
# 2 つのマス
# (x,y) と
# (x
# ′
#  ,y
# ′
#  ) が隣り合っているとは、
# ∣x−x
# ′
#  ∣+∣y−y
# ′
#  ∣=1 であることを指します。

# 赤色に塗られたマスを一様ランダムに
# 1 つ選び、緑色に塗り替えたとき、塗り替え後のグリッドの緑の連結成分数の期待値を
# mod 998244353 で出力してください。


from collections import defaultdict
import sys
from typing import Callable, DefaultDict, List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加

    初始化的连通分量个数 为 n
    """

    __slots__ = ("n", "part", "_parent", "_rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self._parent = list(range(n))
        self._rank = [1] * n

    def find(self, x: int) -> int:
        while self._parent[x] != x:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        """按秩合并."""
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

    def unionTo(self, child: int, parent: int) -> bool:
        """定向合并.将child的父节点设置为parent."""
        rootX = self.find(child)
        rootY = self.find(parent)
        if rootX == rootY:
            return False
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


GREEN = "#"
RED = "."
if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    grid = [list(input()) for _ in range(ROW)]
    uf = UnionFindArray(ROW * COL)
    greenCount = 0
    for i in range(ROW):
        for j in range(COL):
            if grid[i][j] == GREEN:
                greenCount += 1
                cur = i * COL + j
                if i - 1 >= 0 and grid[i - 1][j] == GREEN:
                    next = (i - 1) * COL + j
                    uf.union(cur, next)
                if j + 1 < COL and grid[i][j + 1] == GREEN:
                    next = i * COL + j + 1
                    uf.union(cur, next)
    redCount = ROW * COL - greenCount
    curPart = uf.part - redCount

    # 红 -> 绿
    res = 0
    for i in range(ROW):
        for j in range(COL):
            if grid[i][j] == RED:
                cur = i * COL + j
                preRoots = set()
                if i - 1 >= 0 and grid[i - 1][j] == GREEN:
                    preRoots.add(uf.find((i - 1) * COL + j))
                if j + 1 < COL and grid[i][j + 1] == GREEN:
                    preRoots.add(uf.find(i * COL + j + 1))
                if i + 1 < ROW and grid[i + 1][j] == GREEN:
                    preRoots.add(uf.find((i + 1) * COL + j))
                if j - 1 >= 0 and grid[i][j - 1] == GREEN:
                    preRoots.add(uf.find(i * COL + j - 1))
                res += curPart - (len(preRoots) - 1)

    print(res * pow(redCount, MOD - 2, MOD) % MOD)
