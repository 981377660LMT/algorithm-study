# 100076. 无限数组的最短子数组
# https://leetcode.cn/problems/minimum-size-subarray-in-infinite-array/
# !求无限循环数组中和为 target 的最短子数组的长度.不存在则返回 -1.
# 1 <= nums.length <= 1e5
# 1 <= nums[i] <= 1e5
# 1 <= target <= 1e9

# !等价于在nums+nums中找到和为target%sum(nums)的最短子数组.哈希表+前缀和即可.

from collections import defaultdict
from typing import List

INF = int(1e18)


class Solution:
    def minSizeSubarray(self, nums: List[int], target: int) -> int:
        def shortestSubarrayWithSumk(arr: List[int], k: int) -> int:
            """在arr中找到和为k的最短子数组的长度.不存在则返回INF."""
            if k == 0:
                return 0
            preSum = defaultdict(int, {0: -1})  # 如果记录索引就是{0: -1}
            res, curSum = INF, 0
            for i, num in enumerate(arr):
                curSum += num
                if curSum - k in preSum:
                    res = min(res, i - preSum[curSum - k])
                preSum[curSum] = i
            return res

        allSum = sum(nums)
        res = shortestSubarrayWithSumk(nums + nums, target % allSum)
        if res == INF:
            return -1
        return res + (target // allSum) * len(nums)


# [1,6,5,5,1,1,2,5,3,1,5,3,2,4,6,6]
print(Solution().minSizeSubarray([1, 6, 5, 5, 1, 1, 2, 5, 3, 1, 5, 3, 2, 4, 6, 6], 56))
# [1,2,3]
print(Solution().minSizeSubarray([1, 2, 3], 3))
