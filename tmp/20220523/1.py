from math import floor
from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def percentageLetter(self, s: str, letter: str) -> int:
        return floor(s.count(letter) / len(s) * 100)

