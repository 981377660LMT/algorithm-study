class Solution:
    def solve(self, nums):
        nums = sorted(nums)
        return max(nums[-1] * nums[-2] * nums[-3], nums[0] * nums[1] * nums[-1])


# 数组中三个数的最大乘积 结论

