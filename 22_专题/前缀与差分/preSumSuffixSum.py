from typing import Callable, Generic, List, Sequence, TypeVar


T = TypeVar("T")


class PreSumSuffixSum(Generic[T]):
    __slots__ = ("_e", "_op", "_preSum", "_suffixSum")

    def __init__(self, seq: Sequence[T], e: Callable[[], T], op: Callable[[T, T], T]) -> None:
        self._e = e
        self._op = op
        n = len(seq)
        preSum: List[T] = [None] * (n + 1)  # type: ignore
        suffixSum: List[T] = [None] * (n + 1)  # type: ignore
        preSum[0] = e()
        suffixSum[n] = e()
        for i in range(n):
            preSum[i + 1] = op(preSum[i], seq[i])
            suffixSum[n - i - 1] = op(suffixSum[n - i], seq[n - i - 1])
        self._preSum = preSum
        self._suffixSum = suffixSum

    def preSum(self, end: int) -> T:
        """查询前缀[0,end)的和."""
        if end < 0:
            return self._e()
        if end >= len(self._preSum):
            return self._preSum[-1]
        return self._preSum[end]

    def suffixSum(self, start: int) -> T:
        """查询后缀[start,n)的和."""
        if start < 0:
            return self._suffixSum[0]
        if start >= len(self._suffixSum):
            return self._e()
        return self._suffixSum[start]


if __name__ == "__main__":
    # 100086. 有序三元组中的最大值 II
    # https://leetcode.cn/problems/maximum-value-of-an-ordered-triplet-ii/
    class Solution:
        def maximumTripletValue(self, nums: List[int]) -> int:
            P = PreSumSuffixSum(nums, lambda: 0, lambda a, b: a if a > b else b)
            res = 0
            for j in range(1, len(nums) - 1):
                preMax = P.preSum(j)
                sufMax = P.suffixSum(j + 1)
                res = max(res, (preMax - nums[j]) * sufMax)
            return res
