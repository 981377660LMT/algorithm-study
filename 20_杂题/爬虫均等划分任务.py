from typing import List, TypeVar
from itertools import accumulate

T = TypeVar('T')


def divide(nums: List[T], k: int) -> List[List[T]]:
    n = len(nums)
    base, plus = divmod(n, k)
    groupLen = [base + int(i < plus) for i in range(k)]
    groupStart = [0] + list(accumulate(groupLen))
    return [nums[start:stop] for start, stop in zip(groupStart, groupStart[1:])]


print(divide(list(range(10)), 4))

