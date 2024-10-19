import bisect
from collections import defaultdict
from heapq import heappop, heappush
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def rearrange(A, order: List[int]) -> List[int]:
    """将数组按照order里的顺序重新排序."""
    res = [None] * len(A)
    for i in range(len(order)):
        res[i] = A[order[i]]
    return res


def argSort(A):
    """返回排序后的索引."""
    return sorted(range(len(A)), key=lambda i: A[i])


class TopKSum:
    """
    默认是`最小`的k个数之和.
    """

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
    T = int(input())
    for _ in range(T):
        N, K = map(int, input().split())
        A = list(map(int, input().split()))
        B = list(map(int, input().split()))

        orderA = argSort(A)
        sortedA = rearrange(A, orderA)
        sortedB = rearrange(B, orderA)

        res = INF
        S = TopKSum(K)

        res, curSum = INF, 0
        for right in range(N):
            S.add(sortedB[right])

            if right >= K - 1:
                res = min(res, sortedA[right] * S.query())
        print(res)
