from itertools import groupby
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个二进制数组 nums 。

# 如果一个子数组中 不存在 两个 相邻 元素的值 相同 的情况，我们称这样的子数组为 交替子数组 。


# 返回数组 nums 中交替子数组的数量
class Solution:
    def countAlternatingSubarrays(self, nums: List[int]) -> int:
        arr = [v ^ (i & 1) for i, v in enumerate(nums)]
        groups = [(char, len(list(group))) for char, group in groupby(arr)]
        res = 0
        for _, length in groups:
            res += length * (length + 1) // 2
        return res


# nums = [0,1,1,1]

print(Solution().countAlternatingSubarrays([0, 1, 1, 1]))
