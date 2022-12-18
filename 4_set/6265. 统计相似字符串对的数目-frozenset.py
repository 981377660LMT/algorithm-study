from typing import List
from collections import Counter


# 给你一个下标从 0 开始的字符串数组 words 。

# 如果两个字符串由相同的字符组成，则认为这两个字符串 相似 。

# 例如，"abca" 和 "cba" 相似，因为它们都由字符 'a'、'b'、'c' 组成。
# 然而，"abacba" 和 "bcfd" 不相似，因为它们不是相同字符组成的。
# 请你找出满足字符串 words[i] 和 words[j] 相似的下标对 (i, j) ，并返回下标对的数目，其中 0 <= i < j <= word.length - 1 。


class Solution:
    def similarPairs(self, words: List[str]) -> int:
        counter = Counter()
        for word in words:
            counter[frozenset(word)] += 1
        return sum(v * (v - 1) // 2 for v in counter.values())
