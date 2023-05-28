from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


# 给你一个用字符串表示的正整数 num ，请你以字符串形式返回不含尾随零的整数 num 。
class Solution:
    def removeTrailingZeros(self, num: str) -> str:
        return num.rstrip("0")
