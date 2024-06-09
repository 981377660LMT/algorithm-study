from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s 。

# 你的任务是重复以下操作删除 所有 数字字符：


# 删除 第一个数字字符 以及它左边 最近 的 非数字 字符。
# 请你返回删除所有数字字符以后剩下的字符串。
class Solution:
    def clearDigits(self, s: str) -> str:
        while any(c.isdigit() for c in s):
            for i in range(len(s)):
                if s[i].isdigit():
                    j = i - 1
                    while j >= 0 and s[j].isdigit():
                        j -= 1
                    if j == -1:
                        s = s[i + 1 :]
                    else:
                        s = s[:j] + s[i + 1 :]
                    break
        return s


print(Solution().clearDigits("cb34"))
