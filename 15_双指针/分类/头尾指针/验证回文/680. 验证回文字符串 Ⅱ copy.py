# 680. 验证回文串 II
# https://leetcode.cn/problems/valid-palindrome-ii/description/


class Solution:
    def validPalindrome(self, s: str) -> bool:
        """给定一个非空字符串 s, 最多`删除`一个字符。判断是否能成为回文字符串。"""

        def isPalindrome(t: str) -> bool:
            return t == t[::-1]

        left, right = 0, len(s) - 1
        while left < right:
            if s[left] != s[right]:
                return isPalindrome(s[left:right]) or isPalindrome(s[left + 1 : right + 1])
            left += 1
            right -= 1
        return True
