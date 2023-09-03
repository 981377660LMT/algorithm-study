# https://leetcode.cn/problems/find-the-minimum-possible-sum-of-a-beautiful-array/
# 2834. 找出美丽数组的最小和 (k-void数组)

# 给你两个正整数：n 和 target 。
# 如果数组 nums 满足下述条件，则称其为 美丽数组 。
# nums.length == n.
# nums 由两两互不相同的正整数组成。
# !在范围 [0, n-1] 内，不存在 两个 不同 下标 i 和 j ，使得 nums[i] + nums[j] == target 。
# 返回符合条件的美丽数组所可能具备的 最小 和。


class Solution:
    def minimumPossibleSum(self, n: int, target: int) -> int:
        def getSum(first: int, diff: int, count: int) -> int:
            """等差数列求和"""
            last = first + (count - 1) * diff
            return (first + last) * count // 2

        half = target // 2
        if half >= n:
            return getSum(1, 1, n)

        return getSum(1, 1, half) + getSum(target, 1, n - half)
