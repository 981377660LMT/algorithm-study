# RectangleSum

from bisect import bisect_left
from typing import List, Tuple


class CompressedPointAddRectangleSum:
    __slots__ = ("xs", "ys", "idxs", "comp", "mat")

    def __init__(self, points: List[Tuple[int, int, int]]):
        """二维矩形区域计数 预先添加可能出现的点

        Args:
            points (List[Tuple[int, int, int]]): x, y, weight
        """
        points = sorted(points)
        self.xs, ys, ws = zip(*points)
        self.ys = sorted(set(ys))
        self.idxs = {}
        for i, (x, y, _) in enumerate(points):
            self.idxs[x, y] = i
        self.comp = {val: idx for idx, val in enumerate(self.ys)}
        ys = [self.comp[val] for val in ys]
        MAXLOG = len(self.ys).bit_length()
        self.mat = _PointAddRectangleSum(ys, ws, MAXLOG)

    def add(self, x: int, y: int, val: int):
        idx = self.idxs[x, y]
        self.mat.point_add(idx, val)

    def query(self, x1: int, x2: int, y1: int, y2: int) -> int:
        """求矩形x1<=x<x2,y1<=y<y2的权值和 注意是左闭右开"""
        return self._rect_sum(x1, x2, y2) - self._rect_sum(x1, x2, y1)

    def _rect_sum(self, l, r, upper):
        l = bisect_left(self.xs, l)
        r = bisect_left(self.xs, r)
        upper = bisect_left(self.ys, upper)
        return self.mat.rect_sum(l, r, upper)


class _BitVector:
    __slots__ = ("BLOCK_NUM", "bit", "count")

    def __init__(self, size):
        self.BLOCK_NUM = (size + 31) >> 5
        self.bit = [0] * self.BLOCK_NUM
        self.count = [0] * self.BLOCK_NUM

    def __getitem__(self, i):
        return (self.bit[i >> 5] >> (i & 31)) & 1

    def set(self, i):
        self.bit[i >> 5] |= 1 << (i & 31)

    def build(self):
        for i in range(self.BLOCK_NUM - 1):
            self.count[i + 1] = self.count[i] + self.popcount(self.bit[i])

    def popcount(self, x):
        x = x - ((x >> 1) & 0x55555555)
        x = (x & 0x33333333) + ((x >> 2) & 0x33333333)
        x = (x + (x >> 4)) & 0x0F0F0F0F
        x = x + (x >> 8)
        x = x + (x >> 16)
        return x & 0x0000007F

    def rank1(self, r):
        mask = (1 << (r & 31)) - 1
        return self.count[r >> 5] + self.popcount(self.bit[r >> 5] & mask)

    def rank0(self, r):
        mask = (1 << (r & 31)) - 1
        return r - (self.count[r >> 5] + self.popcount(self.bit[r >> 5] & mask))


class _BinaryIndexedTree:
    __slots__ = ("size", "bit")

    def __init__(self, n):
        self.size = n
        self.bit = [0] * (n + 1)

    def build(self, array):
        for i, val in enumerate(array):
            self.bit[i + 1] = val
        for i in range(1, self.size):
            if i + (i & -i) > self.size:
                continue
            self.bit[i + (i & -i)] += self.bit[i]

    def _sum(self, i):
        s = 0
        while i > 0:
            s += self.bit[i]
            i -= i & -i
        return s

    def add(self, i, val):
        i += 1
        while i <= self.size:
            self.bit[i] += val
            i += i & -i

    def sum(self, l, r):
        return self._sum(r) - self._sum(l)


class _PointAddRectangleSum:
    __slots__ = ("n", "array", "mat", "zs", "data", "MAXLOG")

    def __init__(self, array, ws, MAXLOG=32):
        self.MAXLOG = MAXLOG
        self.n = len(array)
        self.array = array
        self.mat = []
        self.zs = []
        self.data = [None for _ in range(self.MAXLOG)]

        order = [i for i in range(self.n)]
        for d in reversed(range(self.MAXLOG)):
            vec = _BitVector(self.n + 1)
            ls = []
            rs = []
            for i, val in enumerate(order):
                if array[val] & (1 << d):
                    rs.append(val)
                    vec.set(i)
                else:
                    ls.append(val)
            vec.build()
            self.mat.append(vec)
            self.zs.append(len(ls))
            order = ls + rs
            self.data[-d - 1] = _BinaryIndexedTree(self.n)
            self.data[-d - 1].build([ws[val] for val in order])

    def point_add(self, k, val):
        y = self.array[k]
        for d in range(self.MAXLOG):
            if y >> (self.MAXLOG - d - 1) & 1:
                k = self.mat[d].rank1(k) + self.zs[d]
            else:
                k = self.mat[d].rank0(k)
            self.data[d].add(k, val)

    def rect_sum(self, l, r, upper):
        res = 0
        for d in range(self.MAXLOG):
            if upper >> (self.MAXLOG - d - 1) & 1:
                res += self.data[d].sum(self.mat[d].rank0(l), self.mat[d].rank0(r))
                l = self.mat[d].rank1(l) + self.zs[d]
                r = self.mat[d].rank1(r) + self.zs[d]
            else:
                l = self.mat[d].rank0(l)
                r = self.mat[d].rank0(r)
        return res


# 0 x y w : 向平面上添加一个点(x,y)并且权值为w
# 1 left down right up : 查询矩形 left<=x<right down<=y<up 内的点权和
# n,q<=1e5 0<=xi,yi<=1e9
n, q = map(int, input().split())
points = [tuple(map(int, input().split())) for _ in range(n)]
queries = [list(map(int, input().split())) for _ in range(q)]

for op, *args in queries:
    if op == 0:
        x, y, _ = args
        points.append((x, y, 0))  # 0: 预先添加可能出现的点

res = []
rectangleSum = CompressedPointAddRectangleSum(points)
for op, *args in queries:
    if op == 0:
        x, y, w = args
        rectangleSum.add(x, y, w)
    else:
        left, down, right, up = args
        res.append(rectangleSum.query(left, right, down, up))

print(*res, sep="\n")
