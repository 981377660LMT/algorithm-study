from typing import List

# 你可以设计一个对数时间复杂度（即，O(log(n))）的解决方案吗？

# 现有一个按 `升序` 排列的整数数组 nums ，其中每个数字都 `互不相同` 。
# 给你一个整数 k ，请你找出并返回从数组最左边开始的第 k 个缺失数字。


def findKthPositive(arr: List[int], k: int) -> int:
    l, r = 0, len(arr) - 1
    while l <= r:
        mid = (l + r) >> 1
        diff = arr[mid] - (mid + 1)
        if diff >= k:
            r = mid - 1
        else:
            l = mid + 1
    return l + k


class Solution:
    def missingElement(self, nums: List[int], k: int) -> int:
        offset = nums[0] - 1
        nums = [num - offset for num in nums]
        print(nums)
        return findKthPositive(nums, k) + offset


print(Solution().missingElement(nums=[4, 7, 9, 10], k=1))  # 5
print(Solution().missingElement(nums=[4, 7, 9, 10], k=3))  # 8
# 解释：缺失数字有 [5,6,8,...]，因此第三个缺失数字为 8 。
