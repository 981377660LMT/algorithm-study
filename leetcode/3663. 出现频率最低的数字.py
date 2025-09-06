from collections import Counter


class Solution:
    def getLeastFrequentDigit(self, n: int) -> int:
        counter = Counter(str(n))
        minCh = min((c, ch) for ch, c in counter.items())[1]
        return int(minCh)
