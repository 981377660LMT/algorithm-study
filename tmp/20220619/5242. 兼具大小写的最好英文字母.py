from string import ascii_lowercase, ascii_uppercase
import string
from typing import List, Tuple
from collections import defaultdict, Counter

MOD = int(1e9 + 7)
INF = int(1e20)

P = list(zip(string.ascii_lowercase, string.ascii_uppercase))[::-1]


class Solution:
    def greatestLetter(self, s: str) -> str:
        S = set(s)
        return next((b for a, b in P if a in S and b in S), "")

