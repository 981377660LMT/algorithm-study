from typing import List

# 如果一个字符串的每一个字符都在 allowed 中，就称这个字符串是 一致字符串 。
# 请你返回 words 数组中 一致字符串 的数目。
class Solution:
    def countConsistentStrings(self, allowed: str, words: List[str]) -> int:
        s = set(allowed)
        return sum(set(word) <= s for word in words)


# 小于等于是子集的意思
