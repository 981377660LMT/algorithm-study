# 每天决定取前缀还是后缀中的最大值
# 求最后整个矩阵取到的最大值

from itertools import accumulate
from typing import List, Tuple


def getMaxes(row: List[int]) -> Tuple[int, int, int]:
    preSum = list(accumulate(row, initial=0))
    suffixSum = list(accumulate(row[::-1], initial=0))[::-1]
    return sum(row), max(preSum), max(suffixSum)


def solve(matrix: List[List[int]]) -> int:
    """这个角色将从一维公园的左边开始。你有N天的时间来玩游戏"""
    sum_, preMax, suffixMax = getMaxes(matrix[0])
    left, right = preMax, sum_
    for row in matrix[1:]:
        sum_, preMax, suffixMax = getMaxes(row)
        left, right = max(left + preMax, right + sum_), max(right + suffixMax, left + sum_)
    return max(left, right)


T = int(input())
for _ in range(T):
    row, col = map(int, input().split())
    matrix = []
    for _ in range(row):
        matrix.append(list(map(int, input().split())))
    print(solve(matrix))

