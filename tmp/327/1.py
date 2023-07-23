from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串数组 words 和一个字符 separator ，请你按 separator 拆分 words 中的每个字符串。

# 返回一个由拆分后的新字符串组成的字符串数组，不包括空字符串 。

# 注意


# separator 用于决定拆分发生的位置，但它不包含在结果字符串中。
# 拆分可能形成两个以上的字符串。
# 结果字符串必须保持初始相同的先后顺序。
class Solution:
    def splitWordsBySeparator(self, words: List[str], separator: str) -> List[str]:
        res = []
        for word in words:
            res += [w for w in word.split(separator) if w]
        return res
