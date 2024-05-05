# 给定一个数组，每次可以选择两个正数同时减 1，问做多能操作多少次.

from typing import List


def solve(nums: List[int]) -> int:
    if len(nums) <= 1:
        return 0
    max_ = max(nums)
    sum_ = sum(nums)
    restSum = sum_ - max_
    if max_ > restSum:
        return restSum
    else:
        return sum_ // 2
