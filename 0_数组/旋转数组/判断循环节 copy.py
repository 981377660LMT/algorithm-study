class Solution:
    def repeatedSubstringPattern(self, s: str) -> bool:
        return (s + s).find(s, 1) != len(s)


print(Solution().repeatedSubstringPattern('abab'))
print(Solution().repeatedSubstringPattern('aba'))
print(Solution().repeatedSubstringPattern('aab'))
