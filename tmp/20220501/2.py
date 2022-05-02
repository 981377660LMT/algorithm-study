from collections import defaultdict
from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minimumCardPickup(self, cards: List[int]) -> int:
        adjMap = defaultdict(list)
        for i, char in enumerate(cards):
            adjMap[char].append(i)

        res = -1
        for values in adjMap.values():
            if len(values) > 1:
                for pre, cur in zip(values, values[1:]):
                    cand = cur - pre + 1
                    if res == -1 or cand < res:
                        res = cand
        return res

