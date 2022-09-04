# 给出一个长度为n的序列A,要求最大值为X ,最小值为Y的区间的个数
# n<=2e5
# ai<=2e5

# !注意到子数组最值与子数组长度的单调性，考虑滑动窗口

import sys
import os
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def solve(nums: List[int], lower: int, upper: int) -> int:
    """子数组最小值为lower,最大值为upper的子数组个数

    遍历右端点R,看左端点的取值个数,维护两个最值点最右边的位置
    """
    n = len(nums)
    res, left = 0, 0
    pos1, pos2 = -1, -1
    for right in range(n):
        if nums[right] == lower:
            pos1 = right
        if nums[right] == upper:
            pos2 = right
        if nums[right] < lower or nums[right] > upper:
            left = right + 1
        res += max(0, min(pos1, pos2) - left + 1)

    return res


def main() -> None:
    n, max_, min_ = map(int, input().split())
    nums = list(map(int, input().split()))
    print(solve(nums, min_, max_))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
