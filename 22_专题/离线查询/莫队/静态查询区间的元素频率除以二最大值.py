# https://atcoder.jp/contests/abc242/submissions/32860087

from collections import defaultdict
from MoAlgo import MoAlgo


class QueryMaxMoAlgo(MoAlgo[int, int]):
    """静态查询区间 `元素频率//2` 的和 因为每个变化都是±1 所以可以O(1)维护"""

    def __init__(self, n: int, q: int):
        super().__init__(n, q)
        self._pair = 0
        self._counter = defaultdict(int)

    def _add(self, index: int, delta: int) -> None:
        value = nums[index]
        self._pair -= self._counter[value] // 2
        self._counter[value] += 1
        self._pair += self._counter[value] // 2

    def _remove(self, index: int, delta: int) -> None:
        value = nums[index]
        self._pair -= self._counter[value] // 2
        self._counter[value] -= 1
        self._pair += self._counter[value] // 2

    def _query(self, qLeft: int, qRight: int) -> int:
        return self._pair


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = int(1e9 + 7)

    n = int(input())
    q = int(input())
    nums = list(map(int, input().split()))
    Mo = QueryMaxMoAlgo(n, q)

    for _ in range(q):
        left, right = map(int, input().split())
        left, right = left - 1, right - 1
        Mo.addQuery(left, right)

    res = Mo.work()
    print(*res, sep="\n")
