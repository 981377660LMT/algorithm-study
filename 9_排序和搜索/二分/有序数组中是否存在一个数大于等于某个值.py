from bisect import bisect_left
from typing import List
from functools import partial


def queryExist(nums: List[float], query: int) -> bool:
    if not nums:
        return False
    index = bisect_left(nums, query)
    return index != len(nums)  # 等于则query比前面都大


nums = [1, 2, 3, 4, 5]
queryExistInNums = partial(queryExist, nums)
print(queryExistInNums(6))
print(queryExistInNums(3))
print(queryExistInNums(0))
