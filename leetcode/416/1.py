from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串数组 message 和一个字符串数组 bannedWords。

# 如果数组中 至少 存在两个单词与 bannedWords 中的任一单词 完全相同，则该数组被视为 垃圾信息。


# 如果数组 message 是垃圾信息，则返回 true；否则返回 false。
class Solution:
    def reportSpam(self, message: List[str], bannedWords: List[str]) -> bool:
        banset = set(bannedWords)
        c = 0
        for word in message:
            if word in banset:
                c += 1
                if c == 2:
                    return True
        return False
