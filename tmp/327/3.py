from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始、长度为 n 的数组 usageLimits 。

# 你的任务是使用从 0 到 n - 1 的数字创建若干组，并确保每个数字 i 在 所有组 中使用的次数总共不超过 usageLimits[i] 次。此外，还必须满足以下条件：


# 每个组必须由 不同 的数字组成，也就是说，单个组内不能存在重复的数字。
# 每个组（除了第一个）的长度必须 严格大于 前一个组。
# 在满足所有条件的情况下，以整数形式返回可以创建的最大组数。


# 最多的数字要在开头使用

from typing import List, Sequence, Union


class Solution:
    def maxIncreasingGroups(self, usageLimits: List[int]) -> int:
        usageLimits.sort(reverse=True)

        def check(mid: int) -> bool:
            """能否分成mid组"""
            res = 0
            for num in nums:
                res += num // mid
            return res >= k

        left, right = 1, int(1e9) + 10
        ok = False
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
                ok = True
            else:
                right = mid - 1
        return right


# usageLimits = [2,1,2]
