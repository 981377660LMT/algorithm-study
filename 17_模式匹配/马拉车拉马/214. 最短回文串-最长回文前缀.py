# 214. 最短回文串-最长回文前缀
# 给定一个字符串 s，你可以通过在字符串前面添加字符将其转换为回文串。
# 找到并返回可以用这种方式转换的最短回文串。
# 0 <= s.length <= 5e4

# !找到最长的回文前缀，然后把剩下的部分反转拼接到前面即可。

from Manacher import Manacher


class Solution:
    def shortestPalindrome(self, s: str) -> str:
        if not s:
            return ""
        M = Manacher(s)
        max_ = max(M.getLongestEvenStartsAt(0), M.getLongestOddStartsAt(0))
        return s[max_:][::-1] + s


print(Solution().shortestPalindrome(s="aacecaaa"))
print(Solution().shortestPalindrome(s=""))
