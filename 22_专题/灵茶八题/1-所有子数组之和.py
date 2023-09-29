# 所有子数组之和
# 给你一个长为 n 的数组 a，输出它的所有连续子数组的元素和的元素和。
# 索引为i的数贡献次数为(i+1)*(n-i)


from typing import List


def sumOfSubarraySum(nums: List[int]) -> int:
    res = 0
    n = len(nums)
    for i, v in enumerate(nums):
        res += (i + 1) * (n - i) * v
    return res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums = list(map(int, input().split()))
    print(sumOfSubarraySum(nums))
