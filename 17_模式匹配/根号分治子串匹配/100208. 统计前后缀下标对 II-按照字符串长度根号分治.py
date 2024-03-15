# https://leetcode.cn/problems/count-prefix-and-suffix-pairs-ii/
# 当 str1 同时是 str2 的前缀和后缀时，isPrefixAndSuffix(str1, str2) 返回 true，否则返回 false。
# 以整数形式，返回满足 i < j 且 isPrefixAndSuffix(words[i], words[j]) 为 true 的下标对 (i, j) 的 数量 。

# !暴力，只对存在相同长度字符串的前后缀进行检查。
# !令 L 表示字符串的总长度，复杂度 O(LsqrtL).


from typing import List
from collections import defaultdict


class Solution:
    def countPrefixSuffixPairs(self, words: List[str]) -> int:
        res = 0
        counter = defaultdict(int)
        visitedLen = set()  # !最多 sqrtL 种长度
        for w in words:
            for i in range(1, len(w) + 1):
                if i in visitedLen and w[:i] == w[-i:]:  # !一个字符最多检查 sqrtL 次
                    res += counter[w[:i]]
            counter[w] += 1
            visitedLen.add(len(w))
        return res
