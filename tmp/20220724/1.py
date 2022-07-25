from typing import List, Tuple, Optional
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def repeatedCharacter(self, s: str) -> str:
        counter = Counter()
        for char in s:
            counter[char] += 1
            if counter[char] > 1:
                return char
        return ""
