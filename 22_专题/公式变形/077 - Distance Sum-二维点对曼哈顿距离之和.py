# 求二维平面上所有点对的曼哈顿距离之和
# n<=2e5
import sys

from itertools import accumulate
from typing import List

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n = int(input())
px = []
py = []
for _ in range(n):
    x, y = map(int, input().split())
    px.append(x)
    py.append(y)


def dist1d(nums: List[int]) -> int:
    n = len(nums)
    nums = sorted(nums)
    preSum = [0] + list(accumulate(nums))
    res = 0
    for i in range(n):
        res += i * nums[i] - preSum[i]
    return res


print(dist1d(px) + dist1d(py))
