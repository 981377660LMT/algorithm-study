# Segment Add Get Min

# 平面上已存在线段segments
# 0 left right a b 增加一条线段 y=ax+b on [left, right)
# 1 x 查询横坐标为x时的最小值 如果没有线段返回INF
from bisect import bisect_left
from typing import List, Tuple

INF = int(1e20)


def lineAddGetMin(segments: List[Tuple[int, int, int, int]], queries: List[List[int]]) -> List[int]:
    X = []
    for op, *args in queries:
        if op == 1:
            X.append(args[0])

    tree = LiChaoTree(X)
    for l, r, a, b in segments:
        line = (a, b)
        tree.add_seg(line, l, r)

    res = []
    for op, *args in queries:
        if op == 0:
            l, r, a, b = args
            line = (a, b)
            tree.add_seg(line, l, r)
        else:
            x = args[0]
            res.append(tree.get_min(x))  # 如果没有线段返回INF
    return res


class LiChaoTree:
    def __init__(self, xs):
        self.INF = INF
        xs = sorted(list(set(xs)))
        n = len(xs)
        self.size = 1 << (n - 1).bit_length()
        self.comp_xs = {x: ind for ind, x in enumerate(xs)}
        self.xs = xs + [self.INF] * (self.size - n)
        self.data = [None] * (self.size + self.size)

    def update(self, line, k, l, r):
        while True:
            if self.data[k] is None:
                self.data[k] = line
                return

            mid = (l + r) >> 1
            lx = self.xs[l]
            mx = self.xs[mid]
            rx = self.xs[r - 1]
            lu = self.f(line, lx) < self.f(self.data[k], lx)
            mu = self.f(line, mx) < self.f(self.data[k], mx)
            ru = self.f(line, rx) < self.f(self.data[k], rx)

            if lu and ru:
                self.data[k] = line
                return
            if not lu and not ru:
                return
            if mu:
                self.data[k], line = line, self.data[k]
            if lu != mu:
                r, k = mid, k << 1
            else:
                l, k = mid, k << 1 | 1

    def add_line(self, line):
        self.update(line, 1, 0, self.size)

    def add_seg(self, line, l, r):
        l = bisect_left(self.xs, l)
        r = bisect_left(self.xs, r)
        l0, r0 = l + self.size, r + self.size
        size = 1
        while l0 < r0:
            if l0 & 1:
                self.update(line, l0, l, l + size)
                l0 += 1
                l += size
            if r0 & 1:
                r0 -= 1
                r -= size
                self.update(line, r0, r, r + size)
            l0 >>= 1
            r0 >>= 1
            size <<= 1

    def f(self, line, x):
        a, b = line
        return a * x + b

    def get_min(self, x):
        k = self.comp_xs[x] + self.size
        res = self.INF
        while k > 0:
            if self.data[k] is not None:
                res = min(res, self.f(self.data[k], x))
            k >>= 1
        return res


n, q = map(int, input().split())
segments = [tuple(map(int, input().split())) for _ in range(n)]
queries = [list(map(int, input().split())) for _ in range(q)]
res = lineAddGetMin(segments, queries)
for r in res:
    print("INFINITY" if r == INF else r)
