# https://leetcode.cn/problems/count-special-triplets/description/
# !3583. 统计特殊三元组-一次遍历枚举右维护左
# 给你一个整数数组 nums。
#
# 特殊三元组 定义为满足以下条件的下标三元组 (i, j, k)：
#
# 0 <= i < j < k < n，其中 n = nums.length
# nums[i] == nums[j] * 2
# nums[k] == nums[j] * 2
# 返回数组中 特殊三元组 的总数。
#
# 由于答案可能非常大，请返回结果对 109 + 7 取余数后的值。
# !枚举右边那个数，维护左边，可以做到一次遍历


from typing import List
from collections import defaultdict


MOD = int(1e9 + 7)


class Solution:
    def specialTriplets(self, nums: List[int]) -> int:
        counter1 = defaultdict(int)
        counter12 = defaultdict(int)
        res = 0
        for v in nums:
            if v & 1 == 0:
                res += counter12[v // 2]
            counter12[v] += counter1[v * 2]
            counter1[v] += 1
        return res % MOD
