from typing import List
from collections import Counter


# 长度相同 的单词 words
# 找出 s 中恰好可以由 words 中所有单词串联形成的子串的起始位置。
class Solution:
    def findSubstring(self, s: str, words: List[str]) -> List[int]:
        if not s or not words:
            return []
        res = []
        n = len(words)
        word_len = len(words[0])
        window_len = n * word_len
        target = Counter(words)

        # 对每个可能的起始位置
        i = 0
        while i + window_len <= len(s):
            sliced = []
            start = i
            for _ in range(n):
                sliced.append(s[start : start + word_len])
                start += word_len

            if Counter(sliced) == target:
                res.append(i)
            i += 1
        return res
