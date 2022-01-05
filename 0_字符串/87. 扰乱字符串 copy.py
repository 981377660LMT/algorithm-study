import functools

# 1 <= s1.length <= 30
class Solution:
    @functools.lru_cache(None)
    def isScramble(self, s1: str, s2: str) -> bool:
        if s1 == s2:
            return True
        if sorted(s1) != sorted(s2):
            return False
        for i in range(1, len(s1)):
            if self.isScramble(s1[:i], s2[:i]) and self.isScramble(s1[i:], s2[i:]):
                return True
            if self.isScramble(s1[:i], s2[-i:]) and self.isScramble(s1[i:], s2[:-i]):
                return True
        return False


# class Solution:
#     def replaceDigits(self, s: str) -> str:
#         return ''.join([chr(ord(s[i - 1]) + int(v)) if i % 2 else v for i, v in enumerate(s)])

