from typing import List
from collections import defaultdict


class Solution:
    def checkIfExist(self, arr: List[int]) -> bool:
        indexes = defaultdict(list)
        for i, num in enumerate(arr):
            indexes[num].append(i)

        for i, num in enumerate(arr):
            target = num * 2
            lis = indexes[target]
            if len(lis) >= 2 or len(lis) == 1 and lis[0] != i:
                return True

        return False


print(Solution().checkIfExist(arr=[10, 2, 5, 3]))
print(Solution().checkIfExist(arr=[3, 1, 7, 11]))

