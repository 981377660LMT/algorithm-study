from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 如果一个整数 n 在 b 进制下（b 为 2 到 n - 2 之间的所有整数）
# 对应的字符串 全部 都是 回文的 ，
# 那么我们称这个数 n 是 严格回文 的。
class Solution:
    def isStrictlyPalindromic(self, n: int) -> bool:
        ...
