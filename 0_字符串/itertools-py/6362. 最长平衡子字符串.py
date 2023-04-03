# https://leetcode.cn/problems/find-the-longest-balanced-substring-of-a-binary-string/
# 给你一个仅由 0 和 1 组成的二进制字符串 s 。
# 如果子字符串中 所有的 0 都在 1 之前 且其中 0 的数量等于 1 的数量，
# 则认为 s 的这个子字符串是平衡子字符串。请注意，空子字符串也视作平衡子字符串。
# 返回  s 中最长的平衡子字符串长度。

# !01分组,取前后最短长度*2


from itertools import groupby


class Solution:
    def findTheLongestBalancedSubstring(self, s: str) -> int:
        groups = [(char, len(list(group))) for char, group in groupby(s)]
        res, pre = 0, 0
        for c, len_ in groups:
            if c == "1":
                res = max(res, min(pre, len_) * 2)
            pre = len_
        return res
