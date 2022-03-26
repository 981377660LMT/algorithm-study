# 等于转化为不大于相减
from collections import Counter
from typing import List


class Solution:
    def solve(self, words: List[str], k: int) -> int:
        """求恰好包含k个不同word的子数组数"""
        return self.helper(words, k) - self.helper(words, k - 1)

    def helper(self, words: List[str], k: int):
        n = len(words)
        if k == 0:
            return 0
        counter = Counter()
        res = 0
        l = 0
        for r in range(n):
            word = words[r]
            counter[word] += 1
            while len(counter) > k:
                counter[words[l]] -= 1
                if counter[words[l]] == 0:
                    del counter[words[l]]
                l += 1
            res += r - l + 1
        return res
