from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数 num ，返回 num 中能整除 num 的数位的数目。

# 如果满足 nums % val == 0 ，则认为整数 val 可以整除 nums
class Solution:
    def countDigits(self, num: int) -> int:
        res = 0
        for i in str(num):
            if int(i) != 0 and num % int(i) == 0:
                res += 1
        return res
