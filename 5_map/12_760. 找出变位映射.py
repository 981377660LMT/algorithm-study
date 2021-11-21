# 例如，给定

# A = [12, 28, 46, 32, 50]
# B = [50, 12, 32, 46, 28]
#

# 需要返回

# [1, 4, 3, 2, 0]

from typing import List
from collections import defaultdict


class Solution:
    def anagramMappings(self, nums1: List[int], nums2: List[int]) -> List[int]:
        mapping = defaultdict(list)
        for i, v in enumerate(nums2):
            mapping[v].append(i)

        return [mapping[num].pop() for num in nums1]

