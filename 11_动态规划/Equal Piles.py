# 每次把最大值变到第二小的值，求最后所有数相等的操作次数
class Solution:
    def solve(self, nums):
        nums = sorted(nums)
        res = 0
        depth = 0
        for pre, cur in zip(nums, nums[1:]):
            if pre != cur:
                depth += 1
            res += depth
        return res


print(Solution().solve(nums=[4, 8, 2]))
