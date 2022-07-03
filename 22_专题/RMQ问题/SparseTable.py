"""https://www.desgard.com/algo/docs/part2/ch03/1-range-max-query/"""


from math import ceil, floor, log2
from typing import Any, Generic, List, TypeVar

T = TypeVar("T", int, float)


class SparseTable(Generic[T]):
    def __init__(self, nums: List[T]):
        n, upper = len(nums), ceil(log2(len(nums))) + 1
        self._n = n
        self._dp1: List[List[Any]] = [[0] * upper for _ in range(n)]
        self._dp2: List[List[Any]] = [[0] * upper for _ in range(n)]
        for i, num in enumerate(nums):
            self._dp1[i][0] = num
            self._dp2[i][0] = num
        for j in range(1, upper):
            for i in range(n):
                if i + (1 << (j - 1)) >= n:
                    break
                self._dp1[i][j] = max(
                    self._dp1[i][j - 1], self._dp1[i + (1 << (j - 1))][j - 1]
                )
                self._dp2[i][j] = min(
                    self._dp2[i][j - 1], self._dp2[i + (1 << (j - 1))][j - 1]
                )

    def query(self, left: int, right: int, *, ismax=True) -> T:
        """[left,right]区间的最大值"""
        # assert 0 <= left <= right < self._n
        k = floor(log2(right - left + 1))
        if ismax:
            return max(self._dp1[left][k], self._dp1[right - (1 << k) + 1][k])
        else:
            return min(self._dp2[left][k], self._dp2[right - (1 << k) + 1][k])


if __name__ == "__main__":
    nums = list(range(100000))
    st = SparseTable(nums)
    assert st.query(32, 636, ismax=True) == 636
    assert st.query(32, 636, ismax=False) == 32

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
