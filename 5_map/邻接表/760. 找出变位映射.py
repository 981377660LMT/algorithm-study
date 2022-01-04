# 例如，给定

# A = [12, 28, 46, 32, 50]
# B = [50, 12, 32, 46, 28]
#

# 需要返回

# [1, 4, 3, 2, 0]

from typing import List
from collections import defaultdict, deque


class Solution:
    def anagramMappings(self, nums1: List[int], nums2: List[int]) -> List[int]:
        indexes = defaultdict(deque)
        for i, v in enumerate(nums2):
            indexes[v].append(i)

        return [indexes[num].popleft() for num in nums1]


print(Solution().anagramMappings([12, 28, 46, 32, 50], [50, 12, 32, 46, 28]))
# 需要返回
# [1, 4, 3, 2, 0]
