from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的字符串数组 words 和两个整数：left 和 right 。

# 如果字符串以元音字母开头并以元音字母结尾，那么该字符串就是一个 元音字符串 ，其中元音字母是 'a'、'e'、'i'、'o'、'u' 。

# 返回 words[i] 是元音字符串的数目，其中 i 在闭区间 [left, right] 内。
VOWEL = set("aeiou")


class Solution:
    def vowelStrings(self, words: List[str], left: int, right: int) -> int:
        res = 0
        for word in words[left : right + 1]:
            if word[0] in VOWEL and word[-1] in VOWEL:
                res += 1
        return res
