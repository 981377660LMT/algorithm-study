# 每次把最大值变到第二小的值，求最后所有数相等的操作次数
class Solution:
    def solve(self, nums):
        nums = sorted(nums)
        res = 0
        dp = 0
        for pre, cur in zip(nums, nums[1:]):
            if pre != cur:
                dp += 1
            res += dp
        return res


print(Solution().solve(nums=[4, 8, 2]))
