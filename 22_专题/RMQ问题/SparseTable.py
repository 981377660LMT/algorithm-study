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
            step = 1 << (k - 1)
            for i in range(self._n - step * 2 + 1):
                t[i] = p[i] if p[i] > p[i + step] else p[i + step]

    def query(self, start: int, end: int) -> int:
        """[start,end)区间的最大值."""
        k = (end - start).bit_length() - 1
        cand1, cand2 = self._dp[k][start], self._dp[k][end - (1 << k)]
        return cand1 if cand1 > cand2 else cand2


class MinSparseTable:
    """求区间最小值的ST表"""

    __slots__ = "_n", "_h", "_dp"

    def __init__(self, arr: List[int]):
        self._n = len(arr)
        self._h = self._n.bit_length()
        self._dp = [[0] * self._n for _ in range(self._h)]
        self._dp[0] = [a for a in arr]
        for k in range(1, self._h):
            t, p = self._dp[k], self._dp[k - 1]
            step = 1 << (k - 1)
            for i in range(self._n - step * 2 + 1):
                t[i] = p[i] if p[i] < p[i + step] else p[i + step]

    def query(self, start: int, end: int) -> int:
        """[start,end)区间的最小值."""
        k = (end - start).bit_length() - 1
        cand1, cand2 = self._dp[k][start], self._dp[k][end - (1 << k)]
        return cand1 if cand1 < cand2 else cand2


T = TypeVar("T")


class SparseTable(Generic[T]):
    """自定义merger的ST表"""

    __slots__ = "_n", "_h", "_dp", "_e", "_op"

    def __init__(self, arr: List[T], e: Callable[[], T], op: Callable[[T, T], T]):
        self._e = e
        self._op = op
        self._n = len(arr)
        self._h = self._n.bit_length()
        self._dp = [[0] * self._n for _ in range(self._h)]
        self._dp[0] = [a for a in arr]  # type: ignore
        for k in range(1, self._h):
            t, p = self._dp[k], self._dp[k - 1]
            step = 1 << (k - 1)
            for i in range(self._n - step * 2 + 1):
                t[i] = op(p[i], p[i + step])  # type: ignore

    def query(self, start: int, end: int) -> T:
        """[start,end)区间的贡献值."""
        k = (end - start).bit_length() - 1
        return self._op(self._dp[k][start], self._dp[k][end - (1 << k)])  # type: ignore


if __name__ == "__main__":
    INF = int(1e18)
    nums = list(range(10000))
    st1 = MinSparseTable(nums)
    st2 = MaxSparseTable(nums)
    st3 = SparseTable(nums, lambda: INF, min)
    assert st1.query(32, 637) == 32
    assert st2.query(32, 637) == 636
    assert st3.query(32, 637) == 32

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
