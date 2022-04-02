# k次操作可以任意改变k个数的值
# 求最长连续子数组
from collections import defaultdict


class Solution:
    def solve(self, nums, k):
        inedxMap = defaultdict(list)
        for i, num in enumerate(nums):
            inedxMap[num].append(i)
        res = 0
        for row in inedxMap.values():
            left = 0
            for right, num in enumerate(row):
                # 中间有几个不num的数
                while row[right] - row[left] - (right - left) > k:
                    left += 1
                res = max(res, right - left + 1 + k)
        return min(res, len(nums))

