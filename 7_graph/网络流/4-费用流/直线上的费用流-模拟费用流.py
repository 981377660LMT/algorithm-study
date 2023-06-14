# 直线上的最小费用流-模拟费用流
# https://maspypy.github.io/library/flow/min_cost_matching_on_line.hpp


from heapq import heappop, heappush
from typing import List, Optional, Tuple


INF = int(1e18)


def minCostMatchingOnLine(from_: List[int], to: List[int]) -> int:
    """
    给定直线上n个点0,1,...,n-1
    每个位置处有from_[i]个人
    现在要移动这些人,使得每个位置有to[i]个人
    求最小的移动距离
    """
    assert sum(from_) >= sum(to)
    n = len(from_)
    f = SlopeTrick([0] * (n + 1), [])
    for i in range(n):
        diff = to[i] - from_[i]
        f.translate(diff)
        f.clear_right()
        f.add_abs(0)
    return f.get_destructively(0)


class SlopeTrick:
    __slots__ = ("_min_f", "_pq_l", "_pq_r", "add_l", "add_r")

    def __init__(self, left: Optional[List[int]] = None, right: Optional[List[int]] = None):
        self.add_l = 0  # 左侧第一个拐点的位置 -> \_/
        self.add_r = 0  # 右侧第一个拐点的位置 \_/ <-
        self._pq_l = [] if left is None else left  # 大根堆
        self._pq_r = [] if right is None else right  # 小根堆
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

    def get_destructively(self, x: int) -> int:
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

    def merge_destructively(self, st: "SlopeTrick"):
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


if __name__ == "__main__":
    from_ = [1, 2, 3, 4, 5]
    to = [3, 3, 1, 1, 1]
    assert minCostMatchingOnLine(from_, to) == 6
    # https://atcoder.jp/contests/kupc2016/tasks/kupc2016_h
    n = int(input())

    from_ = list(map(int, input().split()))
    to = list(map(int, input().split()))
    print(minCostMatchingOnLine(from_, to))
