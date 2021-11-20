from collections import Counter


class Solution:
    def canConstruct(self, ransomNote: str, magazine: str) -> bool:
        c1, c2 = Counter(ransomNote), Counter(magazine)
        return c1 & c2 == c1


# 输入：ransomNote = "aa", magazine = "aab"
# 输出：true
# 输入：ransomNote = "aa", magazine = "ab"
# 输出：false
