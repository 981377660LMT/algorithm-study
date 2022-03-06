from collections import defaultdict
from functools import cmp_to_key
from typing import List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def sortJumbled(self, mapping: List[int], nums: List[int]) -> List[int]:
        def cmpFunc(num: int) -> int:
            return int(''.join(str(mapping[int(char)]) for char in str(num))))
        return sorted(nums, key=cmpFunc)

    