from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的字符串 s ，将 s 中的元素重新 排列 得到新的字符串 t ，它满足：

# 所有辅音字母都在原来的位置上。更正式的，如果满足 0 <= i < s.length 的下标 i 处的 s[i] 是个辅音字母，那么 t[i] = s[i] 。
# 元音字母都必须以他们的 ASCII 值按 非递减 顺序排列。更正式的，对于满足 0 <= i < j < s.length 的下标 i 和 j  ，如果 s[i] 和 s[j] 都是元音字母，那么 t[i] 的 ASCII 值不能大于 t[j] 的 ASCII 值。
# 请你返回结果字母串。


# 元音字母为 'a' ，'e' ，'i' ，'o' 和 'u' ，它们可能是小写字母也可能是大写字母，辅音字母是除了这 5 个字母以外的所有字母。

VOWEL = set("aeiouAEIOU")


class Solution:
    def sortVowels(self, s: str) -> str:
        sb = [""] * len(s)
        vowel = []
        for i, c in enumerate(s):
            if c in VOWEL:
                vowel.append(c)
        vowel.sort(reverse=True)
        j = 0
        for i, c in enumerate(s):
            if c in VOWEL:
                sb[i] = vowel[j]
                j += 1
            else:
                sb[i] = c
        return "".join(sb)
