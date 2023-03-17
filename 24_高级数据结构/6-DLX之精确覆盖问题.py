# 能否在矩阵中找到一个行的集合，使得这些行中，每一列都有且仅有一个数字 1。
# n,m<=500
# 数据保证矩阵中 1 的数量不超过 5000。 (稀疏矩阵)
# https://www.acwing.com/problem/content/1069/


class Dancinglink:
    def __init__(self, n, debug=False):
        self.n = n
        self.exist = [True] * n
        self._left = [i - 1 for i in range(n)]
        self._right = [i + 1 for i in range(n)]
        self._debug = debug

    def pop(self, k):
        if self._debug:
            assert self.exist[k]
        L = self._left[k]
        R = self._right[k]
        if L != -1:
            if R != self.n:
                self._right[L], self._left[R] = R, L
            else:
                self._right[L] = self.n
        elif R != self.n:
            self._left[R] = -1
        self.exist[k] = False

    def left(self, idx, k=1):
        if self._debug:
            assert self.exist[idx]
        res = idx
        while k:
            res = self._left[res]
            if res == -1:
                break
            k -= 1
        return res

    def right(self, idx, k=1):
        if self._debug:
            assert self.exist[idx]
        res = idx
        while k:
            res = self._right[res]
            if res == self.n:
                break
            k -= 1
        return res
