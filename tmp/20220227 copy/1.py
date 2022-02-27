from typing import List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def prefixCount(self, words: List[str], pref: str) -> int:
        res = 0
        for w in words:
            if w.startswith(pref):
                res += 1
        return res
