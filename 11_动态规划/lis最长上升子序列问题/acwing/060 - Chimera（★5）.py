# 求最长山脉序列的长度，满足山脉中没有两个相邻的相同高度的点，且只有一个峰
# 前后缀LIS，枚举中间点
# n<=3e5

import sys
from typing import List
from bisect import bisect_left, bisect_right


sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def caldp(nums: List[int], isStrict=True) -> List[int]:
    """求每个位置处的LIS长度(包括自身)"""
    if not nums:
        return []
    res = [1] * len(nums)
    LIS = [nums[0]]
    for i in range(1, len(nums)):
        if nums[i] > LIS[-1]:
            LIS.append(nums[i])
            res[i] = len(LIS)
        else:
            pos = bisect_left(LIS, nums[i]) if isStrict else bisect_right(LIS, nums[i])
            LIS[pos] = nums[i]
            res[i] = pos + 1
    return res


n = int(input())
nums = list(map(int, input().split()))

up = caldp(nums)
down = caldp(nums[::-1])[::-1]

res = []
for i in range(n):
    res.append(down[i] + up[i] - 1)
print(max(res, default=0))  # !有max 必加default

