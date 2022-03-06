from collections import defaultdict
from typing import List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def sortJumbled(self, mapping: List[int], nums: List[int]) -> List[int]:
        de = defaultdict(int)
        for num in nums:
            sb = []
            for char in str(num):
                sb.append(str(mapping[int(char)]))
            curNum = int(''.join(sb))
            de[num] = curNum
        return sorted(nums, key=lambda x: de[x])

