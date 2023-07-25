# !维护topK之和 (这里topK是最小值)
# TopKSum

from collections import defaultdict
from heapq import heappop, heappush
from random import randint
from typing import List
from sortedcontainers import SortedList


class TopKSum:
    """
    默认是`最小`的k个数之和.
    """

    # in d_in 是大根堆，out d_out 是小根堆
    __slots__ = ("_sum", "_k", "_in", "_out", "_d_in", "_d_out", "_c", "_min")

    def __init__(self, k: int, min=True) -> None:
        self._k = k
        self._sum = 0
        self._in = []
        self._out = []
        self._d_in = []
        self._d_out = []
        self._min = min
        self._c = defaultdict(int)

    def query(self) -> int:
        return self._sum if self._min else -self._sum

    def add(self, x: int) -> None:
        if not self._min:
            x = -x
        self._c[x] += 1
        heappush(self._in, -x)
        self._sum += x
        self._modify()

    def discard(self, x: int) -> None:
        if not self._min:
            x = -x
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
        if not self._min:
            x = -x
        return self._c[x] > 0


if __name__ == "__main__":
    #  brute force
    k = 5
    ts = TopKSum(k, min=False)
    sl = SortedList(key=lambda x: -x)
    for _ in range(1000):
        # add
        x = randint(0, 100)
        ts.add(x)
        sl.add(x)
        ts.add(x)
        sl.add(x)
        assert ts.query() == sum(sl[:k])

        # setK
        k = randint(1, 10)
        ts.set_k(k)
        assert ts.query() == sum(sl[:k])

        # discard
        x = randint(0, 100)
        ts.discard(x)
        sl.discard(x)
        assert ts.query() == sum(sl[:k])

        assert len(ts) == len(sl)
        assert (x in ts) == (x in sl)

    # 2163. 删除元素后和的最小差值
    # https://leetcode.cn/problems/minimum-difference-in-sums-after-removal-of-elements/
    class Solution:
        def minimumDifference(self, nums: List[int]) -> int:
            # 前面最小n个和后面大n个
            n = len(nums) // 3
            minK, maxK = TopKSum(n, min=True), TopKSum(n, min=False)
            for i in range(n):
                minK.add(nums[i])
            for i in range(n, 3 * n):
                maxK.add(nums[i])
            res = minK.query() + maxK.query()
            for i in range(n, 2 * n):
                minK.add(nums[i])
                maxK.discard(nums[i])
                res = min(res, minK.query() + maxK.query())
            return res
