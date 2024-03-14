# https://leetcode.cn/problems/count-prefix-and-suffix-pairs-ii/
# 当 str1 同时是 str2 的前缀和后缀时，isPrefixAndSuffix(str1, str2) 返回 true，否则返回 false。
# 以整数形式，返回满足 i < j 且 isPrefixAndSuffix(words[i], words[j]) 为 true 的下标对 (i, j) 的 数量 。


from collections import defaultdict
from typing import List


class Solution:
    def countPrefixSuffixPairs2(self, words: List[str]) -> int:
        counter = defaultdict(int)
        res = 0
        for w in words:
            for k, v in counter.items():
                if w.startswith(k) and w.endswith(k):
                    res += v
            counter[w] += 1
        return res
