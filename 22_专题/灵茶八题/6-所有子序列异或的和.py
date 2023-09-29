# 6-所有子序列异或的和
# 二进制拆位
# https://leetcode.cn/problems/sum-of-all-subset-xor-totals/description/

from functools import reduce
from typing import List

MOD = int(1e9 + 7)


def sumOfSubsetXor(nums: List[int]) -> int:
    or_ = reduce(lambda x, y: x | y, nums)
    pow_ = pow(2, len(nums) - 1, MOD)
    return pow_ * or_ % MOD


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums = list(map(int, input().split()))
    print(sumOfSubsetXor(nums))
