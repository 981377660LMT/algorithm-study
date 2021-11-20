from math import gcd


# 此题非常经典
class Solution:
    def gcdOfStrings(self, str1: str, str2: str) -> str:
        if str1 + str2 != str2 + str1:
            return ''
        return str1[: gcd(len(str1), len(str2))]


print(Solution().gcdOfStrings("ABABAB", "ABAB"))

