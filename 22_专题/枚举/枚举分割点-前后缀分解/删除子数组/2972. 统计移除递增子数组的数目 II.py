# 2972. 统计移除递增子数组的数目 II
# https://leetcode.cn/problems/count-the-number-of-incremovable-subarrays-ii/
# 给你一个下标从 0 开始的 正 整数数组 nums 。
# 如果 nums 的一个子数组满足：移除这个子数组后剩余元素 严格递增 ，那么我们称这个子数组为 移除递增 子数组。
# 比方说，[5, 3, 4, 6, 7] 中的 [3, 4] 是一个移除递增子数组，
# 因为移除该子数组后，[5, 3, 4, 6, 7] 变为 [5, 6, 7] ，是严格递增的。
# 请你返回 nums 中 移除递增 子数组的总数目。
# 注意 ，剩余元素为空的数组也视为是递增的。
# 同 1574. 删除最短的子数组使剩余数组有序

from typing import List


class Solution:
    def incremovableSubarrayCount(self, nums: List[int]) -> int:
        n = len(nums)
        i, j = 0, len(nums) - 1
        while i + 1 < n and nums[i] < nums[i + 1]:
            i += 1
        if i == n - 1:  # 全部递增
            return n * (n + 1) // 2
        while j - 1 >= 0 and nums[j] > nums[j - 1]:
            j -= 1

        res = (n - j) + (i + 1)  # 只保留前缀和后缀时的答案
        for v in range(i + 1):
            # !找到最左边的left,使得nums[left] > nums[v]
            left, right = j, n - 1
            ok = False
            while left <= right:
                mid = (left + right) // 2
                if nums[mid] > nums[v]:
                    right = mid - 1
                    ok = True
                else:
                    left = mid + 1
            if ok:
                res += n - left
        return res + 1


# nums = [6,5,7,8]

print(Solution().incremovableSubarrayCount(nums=[6, 5, 7, 8]))
# nums = [8,7,6,6]
print(Solution().incremovableSubarrayCount(nums=[8, 7, 6, 6]))
