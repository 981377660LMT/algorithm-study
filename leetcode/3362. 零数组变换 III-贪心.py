# 3362. 零数组变换 III (带权区间覆盖问题)
# https://leetcode.cn/problems/zero-array-transformation-iii/description/
# 给你一个长度为 n 的整数数组 nums 和一个二维数组 queries ，其中 queries[i] = [li, ri] 。
#
# 每一个 queries[i] 表示对于 nums 的以下操作：
# 将 nums 中下标在范围 [li, ri] 之间的每一个元素 最多 减少 1 。
# 坐标范围内每一个元素减少的值相互 独立 。
#
# 零数组 指的是一个数组里所有元素都等于 0 。
#
# 请你返回 `最多 可以从 queries 中删除多少个元素`，使得 queries 中剩下的元素仍然能将 nums 变为一个 零数组 。
# 如果无法将 nums 变为一个 零数组 ，返回 -1 。
#
# !用最少的区间覆盖，使得所有的数都变为0。
#
# !贪心原则：左端点<=i的所有区间中，右端点越大的区间越好.

from typing import List
from heapq import heappop, heappush


class Solution:
    def maxRemoval(self, nums: List[int], queries: List[List[int]]) -> int:
        n = len(nums)
        queries = sorted(queries)
        pq = []  # 最大堆
        diff = [0] * (n + 1)
        curSum, qi = 0, 0
        for i, v in enumerate(nums):
            curSum += diff[i]
            while qi < len(queries) and queries[qi][0] <= i:
                heappush(pq, -queries[qi][1])
                qi += 1
            while curSum < v and pq and -pq[0] >= i:  # 保证区间覆盖
                curSum += 1  # 取出一个区间
                right = -heappop(pq)
                diff[right + 1] -= 1
            if curSum < v:
                return -1
        return len(pq)
