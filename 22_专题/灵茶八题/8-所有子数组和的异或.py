# 8-所有子数组和的异或
# 借位拆位
# https://atcoder.jp/contests/arc092/tasks/arc092_b
# n<=2e5
# nums[i]<100


# TODO


from typing import List


def xorOfSubarraySum(nums: List[int]) -> int:
    ...


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums = list(map(int, input().split()))
    print(xorOfSubarraySum(nums))
