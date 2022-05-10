from bisect import bisect_right
from collections import defaultdict
from typing import List

# 计算 子数组 nums[li...ri] 中 差绝对值的最小值
# 2 <= nums.length <= 1e5
# 1 <= nums[i] <= 100


class Solution:
    def minDifference(self, nums: List[int], queries: List[List[int]]) -> List[int]:
        adjMap = defaultdict(list)
        for i, num in enumerate(nums):
            adjMap[num].append(i)

        keys = sorted(adjMap.keys())
        res = []
        for left, right in queries:
            minDiff = int(1e20)
            pre = -int(1e20)

            for key in keys:
                indexes = adjMap[key]
                # 检查数字是否在区间出现过
                pos = bisect_right(indexes, right) - 1
                if pos >= 0 and indexes[pos] >= left:
                    minDiff = min(minDiff, key - pre)
                    pre = key

            res.append(minDiff if minDiff != int(1e20) else -1)

        return res


print(Solution().minDifference(nums=[1, 3, 4, 8], queries=[[0, 1], [1, 2], [2, 3], [0, 3]]))

