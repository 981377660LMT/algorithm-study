# 3844. 最长的准回文子字符串
# https://leetcode.cn/problems/longest-almost-palindromic-substring/description/
# 给你一个由小写英文字母组成的字符串 s。
# 如果一个子字符串在删除 恰好 一个字符后变成回文字符串，那么这个子字符串就是 准回文串（almost-palindromic）。
# 返回一个整数，表示字符串 s 中最长的 准回文串 的长度。
# 子字符串是字符串中任意连续的、非空 字符序列。
# 回文串是一个 非空 字符串，正着读和反着读都相同。
#
# 中心扩展法


class Solution:
    def almostPalindromic(self, s: str) -> int:
        n = len(s)

        def expand(left: int, right: int) -> int:
            """[left+1, right) 是回文串."""
            while left >= 0 and right < n and s[left] == s[right]:
                left -= 1
                right += 1
            return right - left - 1

        res = 0
        for i in range(2 * n - 1):
            left, right = i // 2, (i + 1) // 2
            while left >= 0 and right < n and s[left] == s[right]:
                left -= 1
                right += 1
            res = max(res, expand(left - 1, right))
            res = max(res, expand(left, right + 1))
            if res >= n:
                return n
        return res
