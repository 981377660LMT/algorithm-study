from typing import List
from collections import Counter


def findSamePair(arr: List[int]) -> int:
    c = Counter(arr)
    return sum(count * (count - 1) // 2 for _, count in c.items() if count > 1)


print(findSamePair([1, 2, 1, 4, 2, 5, 3, 2]))

