from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def garbageCollection(self, garbage: List[str], travel: List[int]) -> int:
        n = len(garbage)
        glass = [0] * n
        paper = [0] * n
        metal = [0] * n
        res = 0
        for i, word in enumerate(garbage):
            res += len(word)
            for char in word:
                if char == "G":
                    glass[i] += 1
                elif char == "P":
                    paper[i] += 1
                elif char == "M":
                    metal[i] += 1

        preSum = [0] + list(accumulate(travel))
        last1 = next((i for i in range(n - 1, -1, -1) if glass[i] > 0), 0)
        last2 = next((i for i in range(n - 1, -1, -1) if paper[i] > 0), 0)
        last3 = next((i for i in range(n - 1, -1, -1) if metal[i] > 0), 0)
        res += preSum[last1] + preSum[last2] + preSum[last3]
        return res
