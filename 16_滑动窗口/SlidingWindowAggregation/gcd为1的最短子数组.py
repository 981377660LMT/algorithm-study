# gcd为1的最短子数组

from SlidingWindowAggregation import SlidingWindowAggregation

from math import gcd
from typing import List

INF = int(1e20)


def minLen(nums: List[int]) -> int:
    """gcd为1的最短子数组.不存在返回INF."""
    n = len(nums)
    S = SlidingWindowAggregation(lambda: 0, gcd)
    res, n = INF, len(nums)
    for right in range(n):
        S.append(nums[right])
        while S and S.query() == 1:
            res = min(res, len(S))
            S.popleft()
    return res


if __name__ == "__main__":
    # 6392. 使数组所有元素变成 1 的最少操作次数
    # https://leetcode.cn/problems/minimum-number-of-operations-to-make-all-array-elements-equal-to-1/
    # 给你一个下标从 0 开始的 正 整数数组 nums 。你可以对数组执行以下操作 任意 次：
    # !选择一个满足 0 <= i < n - 1 的下标 i ，将 nums[i] 或者 nums[i+1] 两者之一替换成它们的最大公约数。
    # !请你返回使数组 nums 中所有元素都等于 1 的 最少 操作次数。如果无法让数组全部变成 1 ，请你返回 -1 。
    class Solution:
        def minOperations(self, nums: List[int]) -> int:
            if gcd(*nums) != 1:
                return -1
            ones = nums.count(1)
            if ones:
                return len(nums) - ones
            return minLen(nums) - 1 + len(nums) - 1
