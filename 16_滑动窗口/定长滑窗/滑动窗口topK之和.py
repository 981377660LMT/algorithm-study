# !滑动窗口topK之和 (这里topK是最小值)

from collections import defaultdict
from heapq import heappop, heappush
from typing import List


def windowTopKSum(nums: List[int], windowSize: int, k: int) -> List[int]:
    n = len(nums)
    ts = TopKSum(k)
    res = []
    for right in range(n):
        ts.add(nums[right])
        if right >= windowSize:
            ts.discard(nums[right - windowSize])
        if right >= windowSize - 1:
            res.append(ts.query())
    return res


class TopKSum:
    """默认是`最小`的k个数之和,如果要`最大`的k个数之和,需要手动添加负号"""

    # in d_in 是大根堆，out d_out 是小根堆
    __slots__ = ("_sum", "_k", "_in", "_out", "_d_in", "_d_out", "_c")

    def __init__(self, k: int) -> None:
        self._k = k
        self._sum = 0
        self._in = []
        self._out = []
        self._d_in = []
        self._d_out = []
        self._c = defaultdict(int)

    def query(self) -> int:
        return self._sum

    def add(self, x: int) -> None:
        self._c[x] += 1
        heappush(self._in, -x)
        self._sum += x
        self._modify()

    def discard(self, x: int) -> None:
        if self._c[x] == 0:
            return
        self._c[x] -= 1
        if self._in and -self._in[0] == x:
            self._sum -= x
            heappop(self._in)
        elif self._in and -self._in[0] > x:
            self._sum -= x
            heappush(self._d_in, -x)
        else:
            heappush(self._d_out, x)
        self._modify()

    def set_k(self, k: int) -> None:
        self._k = k
        self._modify()

    def get_k(self) -> int:
        return self._k

    def _modify(self) -> None:
        while self._out and (len(self._in) - len(self._d_in) < self._k):
            p = heappop(self._out)
            if self._d_out and p == self._d_out[0]:
                heappop(self._d_out)
            else:
                self._sum += p
                heappush(self._in, -p)

        while len(self._in) - len(self._d_in) > self._k:
            p = -heappop(self._in)
            if self._d_in and p == -self._d_in[0]:
                heappop(self._d_in)
            else:
                self._sum -= p
                heappush(self._out, p)
        while self._d_in and self._in[0] == self._d_in[0]:
            heappop(self._in)
            heappop(self._d_in)

    def __len__(self) -> int:
        return len(self._in) + len(self._out) - len(self._d_in) - len(self._d_out)

    def __contains__(self, x: int) -> bool:
        return self._c[x] > 0


assert windowTopKSum([3, 1, 4, 1, 5, 9], 4, 3) == [5, 6, 10]
assert windowTopKSum([12, 2, 17, 11, 19, 8, 4, 3, 6, 20], 6, 3) == [21, 14, 15, 13, 13]

topKSum = TopKSum(3)
assert topKSum.query() == 0
topKSum.add(3)
assert topKSum.query() == 3
topKSum.add(1)
assert topKSum.query() == 4
topKSum.add(4)
assert topKSum.query() == 8
topKSum.add(1)
assert topKSum.query() == 5
topKSum.discard(3)
assert topKSum.query() == 6
