from typing import List
from collections import Counter


class Solution:
    def findEvenNumbers(self, digits: List[int]) -> List[int]:
        res = []
        store = Counter(map(str, digits))
        for num in range(100, 1000, 2):
            cur = Counter(list(str(num)))
            if store & cur == cur:
                res.append(num)

        return res


print(Solution().findEvenNumbers([2, 1, 3, 0]))
