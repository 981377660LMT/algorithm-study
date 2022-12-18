"""https://www.desgard.com/algo/docs/part2/ch03/1-range-max-query/"""


from typing import Callable, Generic, List, TypeVar


class MaxSparseTable:
    """求区间最大值的ST表"""

    __slots__ = "_n", "_dp"

    def __init__(self, arr: List[int]):
        n = len(arr)
        size = n.bit_length()
        self._n = n

        dp = [[0] * n for _ in range(size)]  # !dp[i][j]表示闭区间[j,j+2**i-1]的最大值
        dp[0] = arr[:]

        for i in range(1, size):
            for j in range(n - (1 << i) + 1):
                cand1, cand2 = dp[i - 1][j], dp[i - 1][j + (1 << (i - 1))]
                dp[i][j] = cand1 if cand1 > cand2 else cand2
        self._dp = dp

    def query(self, left: int, right: int) -> int:
        """[left,right]区间的最大值"""
        assert 0 <= left <= right < self._n, f"{left}, {right}, {self._n}"
        k = (right - left + 1).bit_length() - 1
        cand1, cand2 = self._dp[k][left], self._dp[k][right - (1 << k) + 1]
        if cand1 > cand2:
            return cand1
        return cand2


class MinSparseTable:
    """求区间最小值的ST表"""

    __slots__ = "_n", "_dp"

    def __init__(self, arr: List[int]):
        n = len(arr)
        size = n.bit_length()
        self._n = n

        dp = [[0] * n for _ in range(size)]
        dp[0] = arr[:]

        for i in range(1, size):
            for j in range(n - (1 << i) + 1):
                cand1, cand2 = dp[i - 1][j], dp[i - 1][j + (1 << (i - 1))]
                dp[i][j] = cand1 if cand1 < cand2 else cand2
        self._dp = dp

    def query(self, left: int, right: int) -> int:
        """[left,right]区间的最小值"""
        assert 0 <= left <= right < self._n, f"{left}, {right}, {self._n}"
        k = (right - left + 1).bit_length() - 1
        cand1, cand2 = self._dp[k][left], self._dp[k][right - (1 << k) + 1]
        if cand1 < cand2:
            return cand1
        return cand2


T = TypeVar("T")
Merger = Callable[[T, T], T]


class SparseTable(Generic[T]):
    """自定义merger的ST表"""

    __slots__ = "_n", "_dp", "_merger"

    def __init__(self, arr: List[T], merger: Merger[T]):
        n = len(arr)
        size = n.bit_length()
        self._n = n
        self._merger = merger

        dp = [[0] * n for _ in range(size)]
        dp[0] = arr[:]  # type: ignore

        for i in range(1, size):
            for j in range(n - (1 << i) + 1):
                dp[i][j] = merger(dp[i - 1][j], dp[i - 1][j + (1 << (i - 1))])  # type: ignore
        self._dp = dp

    def query(self, left: int, right: int) -> T:
        """[left,right]区间的贡献值"""
        k = (right - left + 1).bit_length() - 1
        return self._merger(self._dp[k][left], self._dp[k][right - (1 << k) + 1])  # type: ignore


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
