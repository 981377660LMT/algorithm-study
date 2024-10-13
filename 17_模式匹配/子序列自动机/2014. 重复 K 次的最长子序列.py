"""
# 2014. 重复 K 次的最长子序列
# https://leetcode.cn/problems/longest-subsequence-repeated-k-times/
给你一个长度为 n 的字符串 s ，和一个整数 k 。
请你找出字符串 s 中 重复 k 次的 最长子序列 。
如果存在多个满足的子序列，则返回 字典序最大 的那个。
如果不存在这样的子序列，返回一个 空 字符串。
2<=k<=2000
2<=n<k*8

!按照数据,合法的字符不超过7种 => 考虑倒序枚举长度+全排列+check
"""

from collections import Counter
from itertools import permutations
from SubsequenceAutomaton import SubsequenceAutomaton1


class Solution:
    def longestSubsequenceRepeatedK(self, s: str, k: int) -> str:
        counter = Counter(s)
        # !按照数据 合法的字符不超过7种 => 暴力检验  倒序可以提前返回最大字典序
        cands = []
        for char in sorted(counter, reverse=True):
            repeat = counter[char] // k
            cands.append(char * repeat)
        okChars = "".join(cands)

        SA = SubsequenceAutomaton1(s)
        for length in range(len(okChars), 0, -1):
            for perm in permutations(okChars, length):
                cur = "".join(perm)
                hit, _ = SA.match(cur * k)
                if hit == len(cur) * k:
                    return cur
        return ""


# print(Solution().longestSubsequenceRepeatedK(s="letsleetcode", k=2))
# "bbabbabbbbabaababab"
# 3
print(Solution().longestSubsequenceRepeatedK(s="bbabbabbbbabaababab", k=3))
