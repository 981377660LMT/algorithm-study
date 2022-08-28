"""每个元音包含偶数次的最长子字符串"""

VOWEL = set("aeiou")


class Solution:
    def findTheLongestSubstring(self, s: str) -> int:
        res, state = 0, 0
        preState = dict({0: -1})
        for i, char in enumerate(s):
            if char in VOWEL:
                state ^= 1 << (ord(char) - 97)
            if state not in preState:
                preState[state] = i
            else:
                res = max(res, i - preState[state])
        return res
