# E - Water Tank (abc359 E)
# https://atcoder.jp/contests/abc359/tasks/abc359_e
# 每一时刻，0位会多一高度的水，如果该水高度高过柱子，且高过1位的水高度，则该高度的水会跑到1位，
# 同理继续判断1位，该水是否跑到2位。问每一位出现水的最早时刻。
#
# !单调栈dp

from typing import List
from 每个元素作为最值的影响范围 import getRange


def waterTank(height: List[int]) -> List[int]:
    n = len(height)
    ranges = getRange(height, isMax=True, isLeftStrict=True, isRightStrict=False)
    dp = [0] * (n + 1)
    for i, h in enumerate(height):
        leftBiggerIndex = ranges[i][0]
        dp[i + 1] = dp[leftBiggerIndex] + (i - leftBiggerIndex + 1) * h
    return [v + 1 for v in dp[1:]]


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    N = int(input())
    H = list(map(int, input().split()))
    res = waterTank(H)
    print(*res)
