# 是否是先单增后单减交替的摆动数组
class Solution:
    def solve(self, nums):
        return nums[0] < nums[1] and all(nums[i] != nums[i + 1] for i in range(len(nums) - 1))
