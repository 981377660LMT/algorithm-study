# 缺失的第一个正数
# 给你一个`未排序`的整数数组 nums ，请你找出其中没有出现的最小的正整数。
# 请你实现时间复杂度为 O(n) 并且只使用常数级别额外空间的解决方案。
# !原地哈希
# 总的来说思路就是把3丢到2号位，把1丢到0号位,...
# 遍历一次数组把大于等于1的和小于数组大小的值放到原数组对应位置，然后再遍历一次数组查当前下标是否和值对应
from typing import List


class Solution:
    def firstMissingPositive(self, nums: List[int]) -> int:
        n = len(nums)
        for i in range(n):
            while 1 <= nums[i] <= n and nums[i] != nums[nums[i] - 1]:
                nums[nums[i] - 1], nums[i] = nums[i], nums[nums[i] - 1]
                # nums[i], nums[nums[i] - 1] = nums[nums[i] - 1], nums[i]
        return next((i + 1 for i, num in enumerate(nums) if num != i + 1), n + 1)


print(Solution().firstMissingPositive([1, 2, 0]))
print(Solution().firstMissingPositive([3, 4, -1, 1]))
