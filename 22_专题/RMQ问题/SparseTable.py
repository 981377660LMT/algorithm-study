"""https://www.desgard.com/algo/docs/part2/ch03/1-range-max-query/"""


from typing import Callable, Generic, List, TypeVar


class MaxSparseTable:
    """求区间最大值的ST表"""

    __slots__ = "_n", "_h", "_dp"

    def __init__(self, arr: List[int]):
        self._n = len(arr)
        self._h = self._n.bit_length()
        self._dp = [[0] * self._n for _ in range(self._h)]
        self._dp[0] = [a for a in arr]  # type: ignore
        for k in range(1, self._h):
            t, p = self._dp[k], self._dp[k - 1]
            l = 1 << (k - 1)
            for i in range(self._n - l * 2 + 1):
                t[i] = p[i] if p[i] > p[i + l] else p[i + l]

    def query(self, left: int, right: int) -> int:
        """[left,right]区间的最大值"""
        assert 0 <= left <= right < self._n, f"{left}, {right}, {self._n}"
        right += 1
        k = (right - left).bit_length() - 1
        cand1, cand2 = self._dp[k][left], self._dp[k][right - (1 << k)]
        if cand1 > cand2:
            return cand1
        return cand2


class MinSparseTable:
    """求区间最小值的ST表"""

    __slots__ = "_n", "_h", "_dp"

    def __init__(self, arr: List[int]):
        self._n = len(arr)
        self._h = self._n.bit_length()
        self._dp = [[0] * self._n for _ in range(self._h)]
        self._dp[0] = [a for a in arr]  # type: ignore
        for k in range(1, self._h):
            t, p = self._dp[k], self._dp[k - 1]
            l = 1 << (k - 1)
            for i in range(self._n - l * 2 + 1):
                t[i] = p[i] if p[i] < p[i + l] else p[i + l]

    def query(self, left: int, right: int) -> int:
        """[left,right]区间的最小值"""
        assert 0 <= left <= right < self._n, f"{left}, {right}, {self._n}"
        right += 1
        k = (right - left).bit_length() - 1
        cand1, cand2 = self._dp[k][left], self._dp[k][right - (1 << k)]
        if cand1 < cand2:
            return cand1
        return cand2


T = TypeVar("T")
Merger = Callable[[T, T], T]


class SparseTable(Generic[T]):
    """自定义merger的ST表"""

    __slots__ = "_n", "_h", "_dp", "_op"

    def __init__(self, arr: List[T], op: Merger[T]):
        self._op = op
        self._n = len(arr)
        self._h = self._n.bit_length()
        self._dp = [[0] * self._n for _ in range(self._h)]
        self._dp[0] = [a for a in arr]  # type: ignore
        for k in range(1, self._h):
            t, p = self._dp[k], self._dp[k - 1]
            l = 1 << (k - 1)
            for i in range(self._n - l * 2 + 1):
                t[i] = op(p[i], p[i + l])  # type: ignore

    def query(self, left: int, right: int) -> T:
        """[left,right]区间的贡献值"""
        assert 0 <= left <= right < self._n, f"{left}, {right}, {self._n}"
        right += 1
        k = (right - left).bit_length() - 1
        return self._op(self._dp[k][left], self._dp[k][right - (1 << k)])  # type: ignore


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
