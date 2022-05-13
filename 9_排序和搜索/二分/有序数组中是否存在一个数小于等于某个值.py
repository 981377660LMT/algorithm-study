from bisect import bisect_right
from typing import List
from functools import partial


def queryExist(nums: List[int], query: int) -> bool:
    if not nums:
        return False
    index = bisect_right(nums, query)
    return index != 0  # 等于则query比后面都小


nums = [1, 2, 3, 4, 5]
queryExistInNums = partial(queryExist, nums)
print(queryExistInNums(6))
print(queryExistInNums(3))
print(queryExistInNums(1))
print(queryExistInNums(0))
