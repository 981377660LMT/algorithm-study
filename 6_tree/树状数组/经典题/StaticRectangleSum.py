# StaticRectangleSum

from bisect import bisect_left
from typing import List, Tuple


class CompressedRectangleSum:
    def __init__(self, points: List[Tuple[int, int, int]]):
        """静态二维矩形区域计数 预先添加可能出现的点

        Args:
            points (List[Tuple[int, int, int]]): x, y, weight
        """
        points = sorted(points)
        self.xs, ys, ws = zip(*points)
        self.ys = sorted(set(ys))
        self.comp = {val: idx for idx, val in enumerate(self.ys)}
        ys = [self.comp[val] for val in ys]
        MAXLOG = len(self.ys).bit_length()
        self.mat = _RectangleSum(ys, ws, MAXLOG)

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


class _RectangleSum:
    __slots__ = ("MAXLOG", "n", "mat", "zs", "data")

    def __init__(self, array, ws, MAXLOG=32):
        self.MAXLOG = MAXLOG
        self.n = len(array)
        self.mat = []
        self.zs = []
        self.data = [[0] * (self.n + 1) for _ in range(self.MAXLOG)]

        order = list(range(self.n))
        for d in range(self.MAXLOG - 1, -1, -1):
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
            for i, val in enumerate(order):
                self.data[-d - 1][i + 1] = self.data[-d - 1][i] + ws[val]

    def rect_sum(self, l, r, upper):
        res = 0
        for d in range(self.MAXLOG):
            if upper >> (self.MAXLOG - d - 1) & 1:
                res += self.data[d][self.mat[d].rank0(r)]
                res -= self.data[d][self.mat[d].rank0(l)]
                l = self.mat[d].rank1(l) + self.zs[d]
                r = self.mat[d].rank1(r) + self.zs[d]
            else:
                l = self.mat[d].rank0(l)
                r = self.mat[d].rank0(r)
        return res


# 平面上有n个点,第i个点的坐标为(xi,yi),权值为wi,现在要进行q次询问
# left down right up => 求矩形 left<=x<right, down<=y<up 的权值和

n, q = map(int, input().split())
points = [tuple(map(int, input().split())) for _ in range(n)]
queries = [list(map(int, input().split())) for _ in range(q)]


res = []
rectangleSum = CompressedRectangleSum(points)
for left, down, right, up in queries:
    res.append(rectangleSum.query(left, right, down, up))

print(*res, sep="\n")
