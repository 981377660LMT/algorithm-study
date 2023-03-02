# 给你一个下标从 0 开始的整数数组 nums 。
# 一开始，所有下标都没有被标记。你可以执行以下操作任意次：
# !选择两个 互不相同且未标记 的下标 i 和 j ，满足 2 * nums[i] <= nums[j] ，标记下标 i 和 j 。
# 请你执行上述操作任意次，返回 nums 中最多可以标记的下标数目。
# n<=1e5 nums[i]<=1e9

# 1. 贪心+双指针 => 右指针从中间开始移动
# 2. 二分答案 => 判断头部尾部是否满足条件

from typing import List


class Solution:
    def maxNumOfMarkedIndices(self, nums: List[int]) -> int:
        """binary search"""

        def check(mid: int) -> bool:
            pre, suf = nums[:mid], nums[-mid:]
            return all(2 * pre[i] <= suf[i] for i in range(mid))

        nums.sort()
        n = len(nums)
        left, right = 1, n // 2
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right * 2

    def maxNumOfMarkedIndices2(self, nums: List[int]) -> int:
        """two pointers"""
        nums.sort()
        n = len(nums)
        pair = 0
        right = n // 2
        visited = [False] * n
        for left in range(n):
            if visited[left]:
                continue
            while right < n and (nums[left] * 2 > nums[right] or visited[right]):
                right += 1
            if right < n and nums[left] * 2 <= nums[right] and not visited[right]:
                visited[left] = visited[right] = True
                pair += 1
        return pair * 2


assert Solution().maxNumOfMarkedIndices([9, 2, 5, 4]) == 4
assert Solution().maxNumOfMarkedIndices2([9, 2, 5, 4]) == 4
