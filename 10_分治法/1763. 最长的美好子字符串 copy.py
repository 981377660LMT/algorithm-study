# 复杂度26n
class Solution:
    def longestNiceSubstring(self, s: str) -> str:
        charSet = set(s)
        for char in s:
            if char.swapcase() not in charSet:
                return max((self.longestNiceSubstring(chunk) for chunk in s.split(char)), key=len)
        return s


print(Solution().longestNiceSubstring(s="YazaAay"))
