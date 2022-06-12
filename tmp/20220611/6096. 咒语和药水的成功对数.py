from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def successfulPairs(self, S: List[int], potions: List[int], success: int) -> List[int]:
        spells = sorted([(num, i) for i, num in enumerate(S)])
        potions.sort()
        right = len(potions) - 1
        res = [0] * len(spells)
        for num, index in spells:
            while right >= 0 and potions[right] * num >= success:
                right -= 1
            res[index] = len(potions) - 1 - right

        return res

