# 给你一个长为 n 的数组 a，输出它的所有连续子数组的元素异或和的异或和。


from typing import List


def xorOfSubarrayXor(nums: List[int]) -> int:
    """答案为所有偶数下标nums[i]的异或和."""
    res = 0
    n = len(nums)
    for i, v in enumerate(nums):
        res ^= (((i + 1) * (n - i)) & 1) * v
    return res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums = list(map(int, input().split()))
    print(xorOfSubarrayXor(nums))
