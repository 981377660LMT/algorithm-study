# 给你一个整数数组 nums，请你返回该数组中恰有四个因数的这些整数的各因数之和。
from math import isqrt


def getFactors(n: int) -> List[int]:
    """返回 n 的所有因数"""
    upper = isqrt(n) + 1
    small, big = [], []

    for i in range(1, upper):
        if n % i == 0:
            small.append(i)
            big.append(n // i)

    if small[-1] == big[-1]:
        small.pop()
    return small + big[::-1]


class Solution:
    def sumFourDivisors(self, nums: List[int]) -> int:
        f = [getFactors(num) for num in nums]
        f = [arr for arr in f if len(arr) == 4]
        return sum([sum(arr) for arr in f])


print(Solution().sumFourDivisors(nums=[21, 4, 7]))
# 输出：32
# 解释：
# 21 有 4 个因数：1, 3, 7, 21
# 4 有 3 个因数：1, 2, 4
# 7 有 2 个因数：1, 7
# 答案仅为 21 的所有因数的和。
