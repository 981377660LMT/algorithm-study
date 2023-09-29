# 给你一个长为 n 的数组 a，输出它的所有子集的元素和的元素和。

from typing import List

MOD = int(1e9 + 7)


def sumOfSubsetSum(nums: List[int]) -> int:
    res = 0
    pow_ = pow(2, len(nums) - 1, MOD)
    for num in nums:
        res += num * pow_
        res %= MOD
    return res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums = list(map(int, input().split()))
    print(sumOfSubsetSum(nums))
