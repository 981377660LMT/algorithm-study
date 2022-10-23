# 求带权距离距离之和的最小值/带权中位数
# !即求 ∑ai*|xi-x| 的最小值 => 找到中位数 (这里ai可以理解为count个数)

from typing import List


class Solution:
    def minCost(self, nums: List[int], weights: List[int]) -> int:
        """
        移动第i个元素1个距离的代价为cost[i]
        求使所有元素相等的最小代价

        找中位数
        """
        pair = sorted(zip(nums, weights), key=lambda x: x[0])
        weightMid = sum(weights) // 2
        curWeight = 0
        mid = -1
        for num, weigt in pair:
            curWeight += weigt
            if curWeight > weightMid:  # !注意是严格大于
                mid = num
                break
        return sum(abs(num - mid) * weight for num, weight in pair)


print(Solution().minCost(nums=[1, 3, 5, 2], weights=[2, 3, 1, 14]))
print(Solution().minCost(nums=[1, 2, 3], weights=[1, 1, 1]))
