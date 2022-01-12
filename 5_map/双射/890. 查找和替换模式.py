from typing import List


class Solution:
    def findAndReplacePattern(self, words: List[str], pattern: str) -> List[str]:
        def isIsomorphic(s: str, t: str) -> bool:
            if len(s) != len(t):
                return False
            return len(set(s)) == len(set(t)) == len(set(zip(s, t)))

        return [w for w in words if isIsomorphic(w, pattern)]
