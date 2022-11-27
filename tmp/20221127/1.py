from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个正整数 n ，找出满足下述条件的 中枢整数 x ：

# 1 和 x 之间的所有元素之和等于 x 和 n 之间所有元素之和。
# 返回中枢整数 x 。如果不存在中枢整数，则返回 -1 。题目保证对于给定的输入，至多存在一个中枢整数。


class Solution:
    def pivotInteger(self, n: int) -> int:
        for x in range(1, n + 1):
            if sum(range(1, x + 1)) == sum(range(x, n + 1)):
                return x
        return -1
