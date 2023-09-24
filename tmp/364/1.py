from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 二进制 字符串 s ，其中至少包含一个 '1' 。

# 你必须按某种方式 重新排列 字符串中的位，使得到的二进制数字是可以由该组合生成的 最大二进制奇数 。

# 以字符串形式，表示并返回可以由给定组合生成的最大二进制奇数。


# 注意 返回的结果字符串 可以 含前导零。
class Solution:
    def maximumOddBinaryNumber(self, s: str) -> str:
        n = len(s)
        ones = s.count("1")
        res = "1" * (ones - 1) + (n - ones) * "0" + "1"
        return res
