from typing import List


class Solution:
    def threeConsecutiveOdds(self, arr: List[int]) -> bool:
        return '111' in ''.join(str(v & 1) for v in arr)

