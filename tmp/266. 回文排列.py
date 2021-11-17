from collections import Counter


class Solution:
    def canPermutePalindrome(self, s: str) -> bool:
        counter = Counter(s)
        return len([v for v in counter.values() if v & 1]) <= 1

