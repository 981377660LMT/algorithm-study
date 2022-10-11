from typing import List

# 只要数字左边的数字都比它小，右边的数字都比它大就能保证这个数字被找到
class Solution:
    def binarySearchableNumbers(self, nums: List[int]) -> int:
        """未排序数组中的可被二分搜索的数"""
        n = len(nums)
        left, right = nums[0], nums[-1]
        pre, suf = [0] * n, [0] * n

        # 记录 nums[i] 前面的最大值 和 后面的最小值
        for i in range(n):
            pre[i], suf[n - i - 1] = left, right
            left, right = max(left, nums[i]), min(right, nums[n - i - 1])

        res = 0
        for i in range(n):
            # 当前数字比前面所有数字大，比后面所有数字小
            if nums[i] >= pre[i] and nums[i] <= suf[i]:
                res += 1

        return res
