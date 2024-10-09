# 3312. 查询排序后的最大公约数
# https://leetcode.cn/problems/sorted-gcd-pair-queries/description/
# !求第k大的公约数对gcd(nums[i], nums[j]), kthGcdPair
#
# 莫比乌斯反演(莫反)
# !本质上是容斥原理

from typing import List
from bisect import bisect_right
from itertools import accumulate


class Solution:
    def gcdValues(self, nums: List[int], queries: List[int]) -> List[int]:
        upper = max(nums) + 1
        c1, c2 = [0] * upper, [0] * upper  # !c2[i]表示gcd为i的二元组对数
        for v in nums:
            c1[v] += 1
        for f in range(1, upper):
            for m in range(f, upper, f):
                c2[f] += c1[m]
        for i in range(1, upper):
            c2[i] = c2[i] * (c2[i] - 1) // 2
        for f in range(upper - 1, 0, -1):
            for m in range(2 * f, upper, f):
                c2[f] -= c2[m]

        presum = list(accumulate(c2))
        return [bisect_right(presum, kth) for kth in queries]
