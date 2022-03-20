from collections import Counter
from typing import List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def divideArray(self, nums: List[int]) -> bool:
        counter = Counter(nums)
        return all(counter[i] % 2 == 0 for i in counter)

