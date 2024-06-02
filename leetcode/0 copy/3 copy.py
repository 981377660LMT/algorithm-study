from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s 。它可能包含任意数量的 '*' 字符。你的任务是删除所有的 '*' 字符。

# 当字符串还存在至少一个 '*' 字符时，你可以执行以下操作：


# 删除最左边的 '*' 字符，同时删除该星号字符左边一个字典序 最小 的字符。如果有多个字典序最小的字符，你可以删除它们中的任意一个。
# 请你返回删除所有 '*' 字符以后，剩余字符连接而成的 字典序最小 的字符串。


class Solution:
    def clearStars(self, s: str) -> str:
        n = len(s)
        removed = [False] * len(s)
