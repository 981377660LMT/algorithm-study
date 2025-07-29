# P4151 [WC2011] 最大XOR和路径
# https://www.luogu.com.cn/problem/P4151
# 考虑一个边权为非负整数的`无向连通图`，节点编号为 0 到 N-1.
# 试求出一条从 0 号节点到 N-1 号节点的路径，使得路径上经过的边的权值的 XOR 和最大。
# !路径可以重复经过某些点或边
#
# !将所有环的异或扔进线性基，答案就是0到n-1的路径的权值与线性基的最大异或和.
# 异或最短路/异或最长路

from typing import DefaultDict, Callable, List, Optional, Tuple
from collections import defaultdict


def maxXorPath(n: int, edges: List[Tuple[int, int, int]], start: int, target: int) -> Optional[int]:
    vs = VectorSpace()
    uf = UnionFindArrayWithDistXor(n)
    for u, v, w in edges:
        rootU, rootV = uf.find(u), uf.find(v)
        if rootU != rootV:
            uf.union(u, v, w)
        else:
            cycleXor = uf.dist(u, v) ^ w
            vs.add(cycleXor)

    if not uf.isConnected(start, target):
        return None
    dist = uf.dist(start, target)
    return vs.getMax(dist)


class VectorSpace:
    __slots__ = ("bases", "_max")

    def __init__(self, nums: Optional[List[int]] = None) -> None:
        self.bases = []
        self._max = 0
        if nums is not None:
            for v in nums:
                self.add(v)

    def add(self, v: int) -> bool:
        for e in self.bases:
            if e == 0 or v == 0:
                break
            v = min2(v, v ^ e)
        if v:
            self.bases.append(v)
            if v > self._max:
                self._max = v
            return True
        return False

    def getMax(self, xorVal=0) -> int:
        res = xorVal
        for e in self.bases:
            res = max2(res, res ^ e)
        return res

    def getMin(self, xorVal=0) -> int:
        res = xorVal
        for e in self.bases:
            res = min2(res, res ^ e)
        return res

    def copy(self) -> "VectorSpace":
        res = VectorSpace()
        res.bases = self.bases[:]
        res._max = self._max
        return res

    def _orthogonalSpace(self, maxDim: int) -> "VectorSpace":
        self._normalize()
        m = maxDim
        tmp = [0] * m
        for e in self.bases:
            tmp[e.bit_length() - 1] = e
        tmp = transpose(m, m, tmp, inplace=True)
        res = VectorSpace()
        for j, v in enumerate(tmp):
            if v & (1 << j):
                continue
            res.add(v | (1 << j))
        return res

    def _normalize(self, reverse=True) -> None:
        for j, v in enumerate(self.bases):
            for i in range(j):
                self.bases[i] = min2(self.bases[i], self.bases[i] ^ v)
        self.bases.sort(reverse=reverse)

    def __len__(self) -> int:
        return len(self.bases)

    def __iter__(self):
        return iter(self.bases)

    def __repr__(self) -> str:
        return repr(self.bases)

    def __contains__(self, v: int) -> bool:
        for e in self.bases:
            if v == 0:
                break
            v = min2(v, v ^ e)
        return v == 0

    def __or__(self, other: "VectorSpace") -> "VectorSpace":
        v1, v2 = self, other
        if len(v1) < len(v2):
            v1, v2 = v2, v1
        res = v1.copy()
        for v in v2.bases:
            res.add(v)
        return res

    def __and__(self, other: "VectorSpace") -> "VectorSpace":
        maxDim = max2(self._max, other._max).bit_length()
        x = self._orthogonalSpace(maxDim)
        y = other._orthogonalSpace(maxDim)
        if len(x) < len(y):
            x, y = y, x
        for v in y.bases:
            x.add(v)
        return x._orthogonalSpace(maxDim)


def transpose(row: int, col: int, matrix1D: List[int], inplace=False) -> List[int]:
    """矩阵转置

    matrix1D:每个元素是状压的数字
    inplace:是否修改原矩阵
    """
    # assert row == len(matrix1D)
    m = matrix1D[:] if not inplace else matrix1D
    log = 0
    max_ = max2(row, col)
    while (1 << log) < max_:
        log += 1
    if len(m) < (1 << log):
        m += [0] * ((1 << log) - len(m))
    w = 1 << log
    mask = 1
    for i in range(log):
        mask = mask | (mask << (1 << i))
    for t in range(log):
        w >>= 1
        mask = mask ^ (mask >> w)
        for i in range(1 << t):
            for j in range(w):
                m[w * 2 * i + j] = ((m[w * (2 * i + 1) + j] << w) & mask) ^ m[w * 2 * i + j]
                m[w * (2 * i + 1) + j] = ((m[w * 2 * i + j] & mask) >> w) ^ m[w * (2 * i + 1) + j]
                m[w * 2 * i + j] = ((m[w * (2 * i + 1) + j] << w) & mask) ^ m[w * 2 * i + j]
    return m[:col]


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class UnionFindArrayWithDistXor:
    """维护到根节点距离的并查集.距离为异或."""

    __slots__ = ("part", "_data", "_potential")

    def __init__(self, n: int):
        self.part = n
        self._data = [-1] * n
        self._potential = [0] * n

    def union(
        self, x: int, y: int, dist: int, cb: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        """
        p(x) = p(y) + dist.
        如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
        """
        dist ^= self.distToRoot(y) ^ self.distToRoot(x)
        x, y = self.find(x), self.find(y)
        if x == y:
            return dist == 0
        if self._data[x] < self._data[y]:
            x, y = y, x
        self._data[y] += self._data[x]
        self._data[x] = y
        self._potential[x] = dist
        self.part -= 1
        if cb is not None:
            cb(y, x)
        return True

    def find(self, x: int) -> int:
        if self._data[x] < 0:
            return x
        r = self.find(self._data[x])
        self._potential[x] ^= self._potential[self._data[x]]
        self._data[x] = r
        return r

    def dist(self, x: int, y: int) -> int:
        """返回x到y的距离`f(x) - f(y)`."""
        return self.distToRoot(x) ^ self.distToRoot(y)

    def distToRoot(self, x: int) -> int:
        """返回x到所在组根节点的距离`f(x) - f(find(x))`."""
        self.find(x)
        return self._potential[x]

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getSize(self, x: int) -> int:
        return -self._data[self.find(x)]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        res = defaultdict(list)
        for i in range(len(self._data)):
            res[self.find(i)].append(i)
        return res


if __name__ == "__main__":

    def p4151():
        n, m = map(int, input().split())
        edges = []
        for _ in range(m):
            u, v, w = map(int, input().split())
            u, v = u - 1, v - 1
            edges.append((u, v, w))

        print(maxXorPath(n, edges, 0, n - 1))
