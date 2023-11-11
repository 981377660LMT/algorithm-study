# 给定一个 0-1 串，可以将一个子段翻转(flip)一次，问最后能达到多少种不同的 1 n<=3e5
# !1的个数肯定是连续变化的 统计这个变化范围
# 1 变为 0 那么1的个数-1
# 0 变为 1 那么1的个数+1
# !把数组的1看为-1 0看为1 即最大子数组和-最小子数组和+1 (1表示原来的)


# !统计个数变为统计子数组和

from typing import List


def cal(nums: List[int], getMax=True) -> int:
    """统计子数组最大最小变化值"""
    res = 0
    dp = 0
    for num in nums:
        if getMax:
            dp = max(dp, 0) + num
            res = max(res, dp)
        else:
            dp = min(dp, 0) + num
            res = min(res, dp)
    return res


def solve(nums: List[int]) -> int:
    nums = [-1 if num == 1 else 1 for num in nums]
    return cal(nums, getMax=True) - cal(nums, getMax=False) + 1


import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n = int(input())
    nums = list(map(int, input().split()))
    print(solve(nums))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
