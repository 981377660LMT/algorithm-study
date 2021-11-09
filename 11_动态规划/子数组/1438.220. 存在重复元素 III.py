import bisect
from typing import List
from sortedcontainers import SortedList

# 给你一个整数数组 nums 和两个整数 k 和 t 。
# 请你判断是否存在 两个不同下标 i 和 j，使得 abs(nums[i] - nums[j]) <= t ，同时又满足 abs(i - j) <= k 。
# 此题SortedList解法类似1438. 绝对差不超过限制的最长连续子数组 copy
# 本题是固定长度求差值，1438 题是固定差值求长度。


class Solution:
    def containsNearbyAlmostDuplicate(self, nums: List[int], k: int, t: int) -> bool:
        sortedList = SortedList()
        left, right, = 0, 0
        while right < len(nums):
            if right - left > k:
                sortedList.remove(nums[left])
                left += 1
            index = bisect.bisect_left(sortedList, nums[right] - t)  # 第一个 >= (num-t)的index
            if index != len(sortedList) and sortedList[index] <= nums[right] + t:
                return True
            sortedList.add(nums[right])
            right += 1
        return False


print(Solution().containsNearbyAlmostDuplicate([1, 2, 3, 1], 3, 0))
print(Solution().containsNearbyAlmostDuplicate([1, 5, 9, 1, 5, 9], 2, 3))
