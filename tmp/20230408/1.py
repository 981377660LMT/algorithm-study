from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的二维整数数组 nums 。

# 返回位于 nums 至少一条 对角线 上的最大 质数 。如果任一对角线上均不存在质数，返回 0 。

# 注意：


# 如果某个整数大于 1 ，且不存在除 1 和自身之外的正整数因子，则认为该整数是一个质数。
# 如果存在整数 i ，使得 nums[i][i] = val 或者 nums[i][nums.length - i - 1]= val ，则认为整数 val 位于 nums 的一条对角线上。
def isPrime(n: int) -> bool:
    if n < 2:
        return False
    for i in range(2, int(n**0.5) + 1):
        if n % i == 0:
            return False
    return True


class Solution:
    def diagonalPrime(self, nums: List[List[int]]) -> int:
        res = 0
        ROW, COL = len(nums), len(nums[0])
        for i in range(ROW):
            for j in range(COL):
                if i == j or i + j == ROW - 1:
                    if isPrime(nums[i][j]):
                        res = max(res, nums[i][j])
        return res
