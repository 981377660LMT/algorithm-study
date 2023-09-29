# 2863. 最长半递减数组
# https://leetcode.cn/problems/maximum-length-of-semi-decreasing-subarrays/description/
# !返回最长的子数组，子数组的第一个元素严格大于最后一个元素
# 即：对每个数，寻找右侧严格小于它的最远位置
# !二维偏序问题需要先对一个维度排序

# 按照值分组, 从小到大遍历, 每次更新最右边界


from collections import defaultdict
from typing import List

INF = int(1e18)


class Solution:
    def maxSubarrayLength(self, nums: List[int]) -> int:
        mp = defaultdict(list)
        for i, num in enumerate(nums):
            mp[num].append(i)

        res = 0
        maxRight = -INF
        for key in sorted(mp):
            pos = mp[key]
            for p in pos:
                res = max(res, maxRight - p + 1)
            maxRight = max(maxRight, max(pos))
        return res
