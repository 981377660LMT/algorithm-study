# 2444. 统计定界子数组的数目
# https://leetcode.cn/problems/count-subarrays-with-fixed-bounds/
# !求最小值等于mink, 最大值等于maxk的子数组个数
# 滑动窗口


from typing import List


class Solution:
    def countSubarrays(self, nums: List[int], minK: int, maxK: int) -> int:
        n = len(nums)
        res, left = 0, 0
        pos1, pos2 = -1, -1  # !合法的边界
        for right in range(n):
            if nums[right] == minK:
                pos1 = right
            if nums[right] == maxK:
                pos2 = right
            if nums[right] < minK or nums[right] > maxK:
                left = right + 1
            res += max(0, min(pos1, pos2) - left + 1)
        return res


assert Solution().countSubarrays(nums=[1, 3, 5, 2, 7, 5], minK=1, maxK=5) == 2
