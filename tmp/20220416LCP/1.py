from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def giveGem(self, gem: List[int], operations: List[List[int]]) -> int:
        for x, y in operations:
            half = gem[x] // 2
            gem[x] = gem[x] - half
            gem[y] = gem[y] + half
        return max(gem) - min(gem)

