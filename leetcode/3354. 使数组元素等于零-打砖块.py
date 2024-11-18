# 3354. 使数组元素等于零
# https://leetcode.cn/problems/make-array-elements-equal-to-zero/description/
# 一维的打砖块游戏
# 一个小球从指定位置出发，数字就是砖块的血量，打到砖块后砖块掉1血小球反弹
# 如果在结束整个过程后, 所有砖块的血量都变成了0, 那么这个选择是有效的。
# 求有效的选择的数量.


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


class Solution:
    def countValidSelections(self, nums: List[int]) -> int:
        S = PreSumSuffixSum(nums, lambda: 0, lambda x, y: x + y)
        res = 0
        for i in range(len(nums)):
            if nums[i] == 0:
                leftSum, rightSum = S.preSum(i), S.suffixSum(i)
                if leftSum == rightSum:
                    res += 2
                elif abs(leftSum - rightSum) <= 1:
                    res += 1
        return res
