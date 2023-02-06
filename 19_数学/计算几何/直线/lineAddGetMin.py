# Line Add Get Min
# 李超线段树


from typing import List, Tuple

# 平面上已存在直线lines
# 0 a b 增加一条直线 y=ax+b
# 1 x 查询横坐标为x时的最小值
def lineAddGetMin(lines: List[Tuple[int, int]], queries: List[List[int]]) -> List[int]:
    X = []
    for op, *args in queries:
        if op == 1:
            X.append(args[0])

    tree = LiChaoTree(X)
    for a, b in lines:
        tree.add_line(a, b)

    res = []
    for op, *args in queries:
        if op == 0:
            tree.add_line(*args)
        else:
            res.append(tree.get_min(*args))
    return res


INF = int(1e18)


def max(x, y):
    if x > y:
        return x
    return y


def min(x, y):
    if x < y:
        return x
    return y


class LiChaoTree:
    def __init__(self, X):
        X = sorted(set(X))
        N = 1 << (len(X) - 1).bit_length()
        self._tree = [None] * (N << 1)
        self._N, self._X = N, X + [INF] * (N - len(X))
        self._X_inv = {x: i for i, x in enumerate(X)}

    def add_line(self, a, b):
        tree, X = self._tree, self._X
        i, l, r = 1, 0, self._N
        while r - l:
            if tree[i] is None:
                tree[i] = (a, b)
                return
            m = (l + r) >> 1
            xl, xm, xr = X[l], X[m], X[r - 1]
            ai, bi = tree[i]
            left = a * xl + b < ai * xl + bi
            mid = a * xm + b < ai * xm + bi
            right = a * xr + b < ai * xr + bi

            if left is right:
                if left:
                    tree[i] = (a, b)
                return
            if mid:
                tree[i], a, b = (a, b), ai, bi
            if left is not mid:
                i, r = i << 1, m
            else:
                i, l = i << 1 | 1, m

    def get_min(self, x):
        i = self._X_inv[x]
        i += self._N
        res = INF
        tree = self._tree
        while i:
            if tree[i] is not None:
                a, b = tree[i]
                res = min(res, a * x + b)
            i >>= 1
        return res


n, q = map(int, input().split())
lines = [tuple(map(int, input().split())) for _ in range(n)]
queries = [list(map(int, input().split())) for _ in range(q)]
print(*lineAddGetMin(lines, queries), sep="\n")
