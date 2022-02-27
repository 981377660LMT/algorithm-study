from collections import Counter
from typing import List, Tuple
from string import ascii_lowercase

MOD = int(1e9 + 7)


class Solution:
    def minSteps(self, s: str, t: str) -> int:
        c1, c2 = Counter(s), Counter(t)
        res = 0
        for char in ascii_lowercase:
            res += abs(c1[char] - c2[char])
        return res

