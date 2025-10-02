"""前缀后缀和."""

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

    class Solution:
        # 100086. 有序三元组中的最大值 II
        # https://leetcode.cn/problems/maximum-value-of-an-ordered-triplet-ii/
        def maximumTripletValue(self, nums: List[int]) -> int:
            P = PreSumSuffixSum(nums, lambda: 0, lambda a, b: a if a > b else b)
            res = 0
            for j in range(1, len(nums) - 1):
                preMax = P.preSum(j)
                sufMax = P.suffixSum(j + 1)
                res = max(res, (preMax - nums[j]) * sufMax)
            return res

        # https://leetcode.cn/problems/find-the-maximum-factor-score-of-array/description/
        def maxScore(self, nums: List[int]) -> int:
            from math import gcd, lcm

            S1 = PreSumSuffixSum(nums, lambda: 0, gcd)
            S2 = PreSumSuffixSum(nums, lambda: 1, lcm)
            res = S1.suffixSum(0) * S2.suffixSum(0)
            for i in range(len(nums)):
                gcd_ = gcd(S1.preSum(i), S1.suffixSum(i + 1))
                lcm_ = lcm(S2.preSum(i), S2.suffixSum(i + 1))
                res = max(res, gcd_ * lcm_)
            return res

        # https://leetcode.cn/problems/distinct-points-reachable-after-substring-removal/description/
        def distinctPoints(self, s: str, k: int) -> int:
            from typing import Tuple

            def e() -> Tuple[int, int]:
                return (0, 0)

            def op(pos1: Tuple[int, int], pos2: Tuple[int, int]) -> Tuple[int, int]:
                return (pos1[0] + pos2[0], pos1[1] + pos2[1])

            mp = {"U": (0, 1), "D": (0, -1), "L": (-1, 0), "R": (1, 0)}
            arr = [mp[c] for c in s]
            P = PreSumSuffixSum(arr, e, op)
            res = set()
            remain = len(s) - k
            for preLen in range(remain + 1):
                sufLen = remain - preLen
                cur = op(P.preSum(preLen), P.suffixSum(len(s) - sufLen))
                res.add(cur)
            return len(res)
