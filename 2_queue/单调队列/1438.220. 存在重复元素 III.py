from typing import List
from sortedcontainers import SortedList


# 给你一个整数数组 nums 和两个整数 k 和 t 。
# 请你判断是否存在 两个不同下标 i 和 j，使得 abs(nums[i] - nums[j]) <= t ，同时又满足 abs(i - j) <= k 。
# 此题SortedList解法类似1438. 绝对差不超过限制的最长连续子数组 copy
# 本题是固定长度求差值，1438 题是固定差值求长度。

# 有序容器维护


class Solution:
    def containsNearbyAlmostDuplicate(self, nums: List[int], k: int, t: int) -> bool:
        sl = SortedList()

        for right in range(len(nums)):
            if right - k - 1 >= 0:
                sl.discard(nums[right - k - 1])

            pos1 = sl.bisect_right(nums[right] - t) - 1
            if pos1 >= 0 and abs(nums[right] - sl[pos1]) <= t:
                return True

            pos2 = sl.bisect_right(nums[right] + t) - 1
            if pos2 >= 0 and abs(nums[right] - sl[pos2]) <= t:
                return True

            sl.add(nums[right])

        return False


print(Solution().containsNearbyAlmostDuplicate([1, 2, 3, 1], 3, 0))
print(Solution().containsNearbyAlmostDuplicate([1, 5, 9, 1, 5, 9], 2, 3))
