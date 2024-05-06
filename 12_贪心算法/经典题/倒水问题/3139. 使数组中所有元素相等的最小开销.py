# 3139. 使数组中所有元素相等的最小开销
# https://leetcode.cn/problems/minimum-cost-to-equalize-array/description/
# 给你一个整数数组 nums 和两个整数 cost1 和 cost2 。你可以执行以下 任一 操作 任意 次：
# 从 nums 中选择下标 i 并且将 nums[i] 增加 1 ，开销为 cost1。
# 选择 nums 中两个 不同 下标 i 和 j ，并且将 nums[i] 和 nums[j] 都 增加 1 ，开销为 cost2 。
# 你的目标是使数组中所有元素都 相等 ，请你返回需要的 最小开销 之和。
# 由于答案可能会很大，请你将它对 109 + 7 取余 后返回。


# 倒水问题，水壶问题
# !验证四个点，f(max_), f(target), f(target+1), f(target-1)

from typing import List


MOD = int(1e9 + 7)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minCostToEqualizeArray(self, nums: List[int], cost1: int, cost2: int) -> int:
        res = self.solve(nums, cost1, cost2)
        return res % MOD

    def solve(self, nums: List[int], cost1: int, cost2: int) -> int:
        nums.sort()
        n = len(nums)
        max_ = nums[-1]
        allSum = sum(nums)
        if len(nums) == 1:
            return 0
        if len(nums) == 2:
            return (nums[1] - nums[0]) * cost1 % MOD
        if cost1 * 2 <= cost2:
            diffSum = sum(max_ - num for num in nums)
            return diffSum * cost1 % MOD

        def cal(target: int) -> int:
            """将所有数变成target的最小开销."""
            if target < max_:
                raise ValueError("invalid target")
            diffSum = target * n - allSum
            maxDiff = target - nums[0]
            otherDiff = diffSum - maxDiff
            if maxDiff <= otherDiff:
                if diffSum & 1 == 0:
                    return (diffSum // 2) * cost2
                return (diffSum // 2) * cost2 + cost1
            return otherDiff * cost2 + (maxDiff - otherDiff) * cost1

        def getFirstTarget() -> int:
            """获取首个合法的target满足 target-nums[0]<= sum(target - num for num in nums[1:])."""

            def check(mid: int) -> bool:
                return mid - nums[0] <= (n - 1) * mid - (allSum - nums[0])

            if check(max_):
                return max_
            left, right = max_, 2 * max_
            while left <= right:
                mid = (left + right) // 2
                if check(mid):
                    right = mid - 1
                else:
                    left = mid + 1
            return left

        res = cal(max_)
        target = getFirstTarget()
        res = min2(res, cal(target))
        res = min2(res, cal(target + 1))
        if target - 1 >= max_:
            res = min2(res, cal(target - 1))
        return res % MOD
