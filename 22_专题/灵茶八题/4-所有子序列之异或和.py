# 给你一个长为 n 的数组 a，输出它的所有子集的元素异或和的异或和。


from typing import List

MOD = int(1e9 + 7)


def xorOfSubsetXor(nums: List[int]) -> int:
    if len(nums) == 1:
        return nums[0]
    return 0


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums = list(map(int, input().split()))
    print(xorOfSubsetXor(nums))
