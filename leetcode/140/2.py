from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个数组 maximumHeight ，其中 maximumHeight[i] 表示第 i 座塔可以达到的 最大 高度。

# 你的任务是给每一座塔分别设置一个高度，使得：


# 第 i 座塔的高度是一个正整数，且不超过 maximumHeight[i] 。
# 所有塔的高度互不相同。
# 请你返回设置完所有塔的高度后，可以达到的 最大 总高度。如果没有合法的设置，返回 -1 。


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maximumTotalSum(self, maximumHeight: List[int]) -> int:
        maximumHeight.sort(reverse=True)

        cur = maximumHeight[0]
        curSum = 0
        for v in maximumHeight:
            cur = min2(cur, v)
            if cur <= 0:
                return -1
            curSum += cur
            cur -= 1
        return curSum


# maximumHeight = [2,3,4,3]
print(Solution().maximumTotalSum([2, 3, 4, 3]))
