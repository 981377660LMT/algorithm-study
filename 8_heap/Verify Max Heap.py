# 验证最大堆
class Solution:
    def solve(self, nums):
        for i in range(len(nums)):
            j = 2 * i + 1
            if j < len(nums) and nums[i] < nums[j]:
                return False
            if j + 1 < len(nums) and nums[i] < nums[j + 1]:
                return False

        return True
