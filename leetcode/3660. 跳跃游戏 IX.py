# 3660. 跳跃游戏 IX
# https://leetcode.cn/problems/jump-game-ix/description/
# 给你一个整数数组 nums。
# 从任意下标 i 出发，你可以根据以下规则跳跃到另一个下标 j：
# 仅当 nums[j] < nums[i] 时，才允许跳跃到下标 j，其中 j > i。
# 仅当 nums[j] > nums[i] 时，才允许跳跃到下标 j，其中 j < i。
# 对于每个下标 i，找出从 i 出发且可以跳跃 任意 次，能够到达 nums 中的 最大值 是多少。
# 返回一个数组 ans，其中 ans[i] 是从下标 i 出发可以到达的最大值。
#
# !每个数可以跳到左边更大的，跳到右边更小的
# 如果当前max<=后缀min，那么形成了墙

from typing import List

INF = int(1e18)


class Solution:
    def maxValue(self, nums: List[int]) -> List[int]:
        n = len(nums)
        sufMin = [INF] * (n + 1)
        for i in range(n - 1, -1, -1):
            sufMin[i] = min(sufMin[i + 1], nums[i])

        res = [0] * n
        ptr = 0
        curMax = -INF
        for i, v in enumerate(nums):
            curMax = max(curMax, v)
            if curMax <= sufMin[i + 1]:
                res[ptr : i + 1] = [curMax] * (i + 1 - ptr)
                ptr = i + 1
                curMax = -INF
        return res
