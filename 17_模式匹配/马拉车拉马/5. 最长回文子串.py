from Manacher import Manacher


class Solution:
    def longestPalindrome(self, s: str) -> str:
        """给你一个字符串 s，找到 s 中最长的回文子串。"""
        manacher = Manacher(s)
        start, maxLen = 0, 0
        for i in range(len(s)):
            len1 = manacher.getLongestOddStartsAt(i)
            if len1 > maxLen:
                maxLen = len1
                start = i

            len2 = manacher.getLongestEvenStartsAt(i)
            if len2 > maxLen:
                maxLen = len2
                start = i

        return s[start : start + maxLen]


print(Solution().longestPalindrome("babad"))
