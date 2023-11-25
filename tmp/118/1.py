from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的字符串数组 words 和一个字符 x 。

# 请你返回一个 下标数组 ，表示下标在数组中对应的单词包含字符 x 。


# 注意 ，返回的数组可以是 任意 顺序。
class Solution:
    def findWordsContaining(self, words: List[str], x: str) -> List[int]:
        return [i for i, word in enumerate(words) if x in word]
