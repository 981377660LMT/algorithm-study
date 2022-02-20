from typing import List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def sumOfThree(self, num: int) -> List[int]:
        if num % 3 != 0:
            return []
        base = num // 3
        return [base - 1, base, base + 1]

