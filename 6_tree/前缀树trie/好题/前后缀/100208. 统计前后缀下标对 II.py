from collections import defaultdict
from typing import List


class Solution:
    def countPrefixSuffixPairs2(self, words: List[str]) -> int:
        counter = defaultdict(int)
        res = 0
        for w in words:
            for k, v in counter.items():
                if w.startswith(k) and w.endswith(k):
                    res += v
            counter[w] += 1
        return res
