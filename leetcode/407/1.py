from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个正整数 n 和 k。

# 你可以选择 n 的 二进制表示 中任意一个值为 1 的位，并将其改为 0。

# 返回使得 n 等于 k 所需要的更改次数。如果无法实现，返回 -1。


class Solution:
    def minChanges(self, n: int, k: int) -> int:
        if n & k == k:
            return n.bit_count() - k.bit_count()
        return -1
