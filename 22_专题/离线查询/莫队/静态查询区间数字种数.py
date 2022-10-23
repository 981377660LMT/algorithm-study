# AcWing 2492. HH的项链

# 静态查询区间数字种数
# 三种做法，莫队，离线树状数组，主席树


from typing import List
from MoAlgo import MoAlgo


class QueryTypeMoAlgo(MoAlgo[int, int]):
    """静态查询区间数字种数"""

    def __init__(self, data: List[int]):
        super().__init__(data)
        self._count = 0
        self._counter = [0] * int(1e6 + 5)

    def _add(self, value: int, index: int, qLeft: int, qRight: int) -> None:
        self._counter[value] += 1
        if self._counter[value] == 1:
            self._count += 1

    def _remove(self, value: int, index: int, qLeft: int, qRight: int) -> None:
        self._counter[value] -= 1
        if self._counter[value] == 0:
            self._count -= 1

    def _query(self) -> int:
        return self._count


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353
    INF = int(4e18)

    n = int(input())
    nums = list(map(int, input().split()))
    q = int(input())
    moAlgo = QueryTypeMoAlgo(nums)
    for _ in range(q):
        left, right = map(int, input().split())
        left, right = left - 1, right - 1
        moAlgo.addQuery(left, right)
    res = moAlgo.work()
    print(*res, sep="\n")
