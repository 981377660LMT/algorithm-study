from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 的整数数组 nums，以及一个正整数 k 。

# 将这个数组划分为一个或多个长度为 3 的子数组，并满足以下条件：


# nums 中的 每个 元素都必须 恰好 存在于某个子数组中。
# 子数组中 任意 两个元素的差必须小于或等于 k 。
# 返回一个 二维数组 ，包含所有的子数组。如果不可能满足条件，就返回一个空数组。如果有多个答案，返回 任意一个 即可。
class Solution:
    def divideArray(self, nums: List[int], k: int) -> List[List[int]]:
        nums.sort()
        n = len(nums)
        div = n // 3

        for i in range(div):
            if nums[3 * i + 2] - nums[3 * i] > k:
                return []
        return [[nums[3 * i], nums[3 * i + 1], nums[3 * i + 2]] for i in range(div)]
