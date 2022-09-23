# 给你一个仅由小写英文字母组成的字符串 s ，返回其 最长字母序连续子字符串 的长度。
from itertools import groupby


class Solution:
    def longestContinuousSubstring(self, s: str) -> int:
        res = 1
        dp = 1
        for i in range(1, len(s)):
            if ord(s[i]) - ord(s[i - 1]) == 1:
                dp += 1
            else:
                dp = 1
            res = max(res, dp)
        return res

    def longestContinuousSubstring2(self, s: str) -> int:
        """给每个位置的字符减去index 如果相等 就说明连续"""
        nums = [ord(s) - index for index, s in enumerate(s)]
        groups = [(char, list(group)) for char, group in groupby(nums)]
        return max(len(group) for _, group in groups)
