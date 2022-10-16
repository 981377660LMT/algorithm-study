from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 非负 整数 num 。如果存在某个 非负 整数 k 满足 k + reverse(k) = num  ，则返回 true ；否则，返回 false 。

# reverse(k) 表示 k 反转每个数位后得到的数字。


class Solution:
    def sumOfNumberAndReverse(self, num: int) -> bool:
        for x in range(10**5 + 1):
            if x + int(str(x)[::-1]) == num:
                return True
        return False
