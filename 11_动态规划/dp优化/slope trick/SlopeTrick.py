# https://ei1333.github.io/library/structure/others/slope-trick.hpp
# !维护上凸一次函数集合(\_/)的数据结构
# !如果函数斜率变化为+-1，那么可以用这个数据结构来维护最小值(最大值)
# 如果不为+-1，那么需要用平衡树维护
# 一般用于高速化dp

# API:
# 初始时 f(x) = 0

# query() -> (最小值, 最小值时的x的最小值, 最小值时的x的最大值)

# add_all(a) -> f(x) += a
# add_a_minus_x(a) -> f(x) += max(a - x, 0)
# add_x_minus_a(a) -> f(x) += max(x - a, 0)
# add_abs(a) -> f(x) += abs(x - a)

# clear_right() -> f(x) = min f(y) (y <= x)
# clear_left() -> f(x) = min f(y) (y >= x)

# shift_range(a, b) -> f(x) = min f(y) (x-b <= y <= x-a) (a <= b)
# translate(a) -> f(x) = f(x - a) 向右平移a.

# get(x) -> 返回f在x处的值. 会破坏f的数据结构.
# merge(g) -> f(x) = f(x) + g(x). 会破坏g的数据结构.

from heapq import heappop, heappush
from typing import Tuple


INF = int(1e18)


class SlopeTrick:
    __slots__ = ("_min_f", "_pq_l", "_pq_r", "add_l", "add_r")

    def __init__(self):
        self.add_l = 0  # 左侧第一个拐点的位置 -> \_/
        self.add_r = 0  # 右侧第一个拐点的位置 \_/ <-
        self._pq_l = []  # 大根堆
        self._pq_r = []  # 小根堆
        self._min_f = 0

    def query(self) -> Tuple[int, int, int]:
        """返回 `f(x)的最小值, f(x)取得最小值时x的最小值和x的最大值`"""
        return self._min_f, self._top_l(), self._top_r()

    def add_all(self, a: int) -> None:
        """f(x) += a"""
        self._min_f += a

    def add_a_minus_x(self, a: int) -> None:
        """
        ```
        add \\__
        f(x) += max(a - x, 0)
        ```
        """
        tmp = a - self._top_r()
        if tmp > 0:
            self._min_f += tmp
        self._push_r(a)
        self._push_l(self._pop_r())

    def add_x_minus_a(self, a: int) -> None:
        """
        ```
        add __/
        f(x) += max(x - a, 0)
        ```
        """
        tmp = self._top_l() - a
        if tmp > 0:
            self._min_f += tmp
        self._push_l(a)
        self._push_r(self._pop_l())

    def add_abs(self, a: int) -> None:
        """
        ```
        add \\/
        f(x) += abs(x - a)
        ```
        """
        self.add_a_minus_x(a)
        self.add_x_minus_a(a)

    def clear_right(self) -> None:
        """
        取前缀最小值.
        ```
        \\/ -> \\_
        f_{new} (x) = min f(y) (y <= x)
        ```
        """
        while self._pq_r:
            self._pq_r.pop()

    def clear_left(self) -> None:
        """
        取后缀最小值.
        ```
        \\/ -> _/
        f_{new} (x) = min f(y) (y >= x)
        ```
        """
        while self._pq_l:
            self._pq_l.pop()

    def shift(self, a: int, b: int) -> None:
        """
        ```
        \\/ -> \\_/
        f_{new} (x) = min f(y) (x-b <= y <= x-a)
        ```
        """
        assert a <= b
        self.add_l += a
        self.add_r += b

    def translate(self, a: int) -> None:
        """
        函数向右平移a
        ```
        \\/. -> .\\/
        f_{new} (x) = f(x - a)
        ```
        """
        self.shift(a, a)

    def get_destructive(self, x: int) -> int:
        """
        y = f(x), f(x) broken
        会破坏f内部左右两边的堆.
        """
        res = self._min_f
        while self._pq_l:
            tmp = self._pop_l() - x
            if tmp > 0:
                res += tmp
        while self._pq_r:
            tmp = x - self._pop_r()
            if tmp > 0:
                res += tmp
        return res

    def merge_destructive(self, st: "SlopeTrick"):
        """
        f(x) += g(x), g(x) broken
        会破坏g(x)的左右两边的堆.
        """
        if len(st) > len(self):
            st._pq_l, self._pq_l = self._pq_l, st._pq_l
            st._pq_r, self._pq_r = self._pq_r, st._pq_r
            st.add_l, self.add_l = self.add_l, st.add_l
            st.add_r, self.add_r = self.add_r, st.add_r
            st._min_f, self._min_f = self._min_f, st._min_f
        while st._pq_r:
            self.add_x_minus_a(st._pop_r())
        while st._pq_l:
            self.add_a_minus_x(st._pop_l())
        self._min_f += st._min_f

    def _push_r(self, a: int) -> None:
        heappush(self._pq_r, a - self.add_r)

    def _top_r(self) -> int:
        if not self._pq_r:
            return INF
        return self._pq_r[0] + self.add_r

    def _pop_r(self) -> int:
        val = self._top_r()
        if self._pq_r:
            heappop(self._pq_r)
        return val

    def _push_l(self, a: int) -> None:
        heappush(self._pq_l, -a + self.add_l)

    def _top_l(self) -> int:
        if not self._pq_l:
            return -INF
        return -self._pq_l[0] + self.add_l

    def _pop_l(self) -> int:
        val = self._top_l()
        if self._pq_l:
            heappop(self._pq_l)
        return val

    def _size(self) -> int:
        return len(self._pq_l) + len(self._pq_r)

    def __len__(self) -> int:
        return self._size()
