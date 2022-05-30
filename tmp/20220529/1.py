from collections import Counter
from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def rearrangeCharacters(self, s: str, target: str) -> int:
        counter = Counter(s)
        counter2 = Counter(target)
        return min(counter[k] // v for k, v in counter2.items())

