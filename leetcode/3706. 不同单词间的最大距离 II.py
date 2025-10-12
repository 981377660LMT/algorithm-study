# 3706. 不同单词间的最大距离 II
# https://leetcode.cn/problems/maximum-distance-between-unequal-words-in-array-ii/description/?envType=problem-list-v2&envId=sZVESpvF
# !最优解的 i 和 j 必然会涉及到数组的首个元素 words[0] 或最后一个元素 words[n-1]。(否则...)

from typing import List


class Solution:
    def maxDistance(self, words: List[str]) -> int:
        n, a, b = len(words), words[0], words[-1]
        for i in range(n):
            if a != words[-i - 1] or b != words[i]:
                return n - i
        return 0
