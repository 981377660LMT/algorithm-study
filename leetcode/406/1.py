from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个仅由数字组成的字符串 s，在最多交换一次 相邻 且具有相同 奇偶性 的数字后，返回可以得到的字典序最小的字符串。


# 如果两个数字都是奇数或都是偶数，则它们具有相同的奇偶性。例如，5 和 9、2 和 4 奇偶性相同，而 6 和 9 奇偶性不同。
class Solution:
    def getSmallestString(self, s: str) -> str:
        cands = [s]
        for i in range(len(s) - 1):
            if int(s[i]) % 2 == int(s[i + 1]) % 2:
                sb = list(s)
                sb[i], sb[i + 1] = sb[i + 1], sb[i]
                cands.append("".join(sb))
        return min(cands)
