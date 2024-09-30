from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 word 和一个 非负 整数 k。


# 返回 word 的 子字符串 中，每个元音字母（'a'、'e'、'i'、'o'、'u'）至少 出现一次，并且 恰好 包含 k 个辅音字母的子字符串的总数。
class Solution:
    def countOfSubstrings(self, word: str, k: int) -> int:
        res = 0
        n = len(word)
        for i in range(n):
            for j in range(i, n):
                s = word[i : j + 1]
                ok1 = all(c in s for c in "aeiou")
                ok2 = sum(c not in "aeiou" for c in s) == k
                res += ok1 and ok2
        return res
