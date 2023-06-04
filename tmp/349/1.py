from itertools import groupby
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的字符串 s ，重复执行下述操作 任意 次：

# 在字符串中选出一个下标 i ，并使 c 为字符串下标 i 处的字符。并在 i 左侧（如果有）和 右侧（如果有）各 删除 一个距离 i 最近 的字符 c 。
# 请你通过执行上述操作任意次，使 s 的长度 最小化 。


# 返回一个表示 最小化 字符串的长度的整数。
class Solution:
    def minimizedStringLength(self, s: str) -> int:
        counter = Counter(s)
        return len(counter)
