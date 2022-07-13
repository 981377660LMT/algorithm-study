"""https://www.desgard.com/algo/docs/part2/ch03/1-range-max-query/"""


from math import ceil, floor, log2
from typing import Callable, Generic, List, TypeVar


class MaxSparseTable:
    """求区间最大值的ST表"""

    __slots__ = "_n", "_dp"

    def __init__(self, arr: List[int]):
        n, upper = len(arr), ceil(log2(len(arr))) + 1
        self._n = n

        dp = [[0] * upper for _ in range(n)]
        for i in range(n):
            dp[i][0] = arr[i]
        for j in range(1, upper):
            for i in range(n):
                if i + (1 << (j - 1)) >= n:
                    break
                cand1, cand2 = dp[i][j - 1], dp[i + (1 << (j - 1))][j - 1]
                if cand1 > cand2:
                    dp[i][j] = cand1
                else:
                    dp[i][j] = cand2
        self._dp = dp

    def query(self, left: int, right: int) -> int:
        """[left,right]区间的最大值"""
        assert 0 <= left <= right < self._n, f"{left}, {right}, {self._n}"
        k = floor(log2(right - left + 1))
        cand1, cand2 = self._dp[left][k], self._dp[right - (1 << k) + 1][k]
        if cand1 > cand2:
            return cand1
        return cand2


class MinSparseTable:
    """求区间最小值的ST表"""

    __slots__ = "_n", "_dp"

    def __init__(self, arr: List[int]):
        n, upper = len(arr), ceil(log2(len(arr))) + 1
        self._n = n

        dp = [[0] * upper for _ in range(n)]
        for i in range(n):
            dp[i][0] = arr[i]
        for j in range(1, upper):
            for i in range(n):
                if i + (1 << (j - 1)) >= n:
                    break
                cand1, cand2 = dp[i][j - 1], dp[i + (1 << (j - 1))][j - 1]
                if cand1 < cand2:
                    dp[i][j] = cand1
                else:
                    dp[i][j] = cand2
        self._dp = dp

    def query(self, left: int, right: int) -> int:
        """[left,right]区间的最小值"""
        assert 0 <= left <= right < self._n, f"{left}, {right}, {self._n}"
        k = floor(log2(right - left + 1))
        cand1, cand2 = self._dp[left][k], self._dp[right - (1 << k) + 1][k]
        if cand1 < cand2:
            return cand1
        return cand2


T = TypeVar("T")
Merger = Callable[[T, T], T]


class SparseTable(Generic[T]):
    """自定义merger的ST表"""

    __slots__ = "_n", "_dp", "_merger"

    def __init__(self, arr: List[T], merger: Merger[T]):
        n, upper = len(arr), ceil(log2(len(arr))) + 1
        self._n = n
        self._merger = merger

        dp: List[List[T]] = [[0] * upper for _ in range(n)]  # type: ignore
        for i in range(n):
            dp[i][0] = arr[i]
        for j in range(1, upper):
            for i in range(n):
                if i + (1 << (j - 1)) >= n:
                    break
                dp[i][j] = merger(dp[i][j - 1], dp[i + (1 << (j - 1))][j - 1])
        self._dp = dp

    def query(self, left: int, right: int) -> T:
        """[left,right]区间的贡献值"""
        assert 0 <= left <= right < self._n, f"{left} {right} {self._n}"
        k = floor(log2(right - left + 1))
        return self._merger(self._dp[left][k], self._dp[right - (1 << k) + 1][k])


if __name__ == "__main__":
    nums = list(range(10000))
    st1 = MinSparseTable(nums)
    st2 = MaxSparseTable(nums)
    st3 = SparseTable(nums, min)
    assert st1.query(32, 636) == 32
    assert st2.query(32, 636) == 636
    assert st3.query(32, 636) == 32

# 通过了 O(1)的方式完成了指定区间任意范围的 RMQ。
# 对于离线海量数据查询的需求完成了最高度的优化。
# 但是由于 ST 算法需要一个 2 倍增的预处理，所以整体的复杂度在 O(nlogn)。
# 如此评估下来，其实如果查询量极少的情况下，
# 我们用暴力法的时间开销 O(n)是优于 ST 算法的，
# 但是 ST 是在大量查询的场景下，所以算法也和业务技术方案一样，
# 有适合于业务的，
# 也有不适合于业务的，一切从业务出发来解决问题就好啦~

# [3, 2, 4, 5, 6, 8, 1, 2, 9, 7]
# f(1,0) 表示第 1 个数起，长度为 2^0 =1 的最大值
# f(1,1)=max(3,2)=3
# f(1, 2) = max(3, 2, 4, 5) = 5
# f(1, 3) = max(3, 2, 4, 5, 6, 8, 1, 2) = 8
