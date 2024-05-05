# !两个区间可以随意移动 求右端点减左端点的最小值 此时是差最大的一个
# !但是如果规定了移动方向就不行了


from typing import List


def solve1(nums: List[int], color: str) -> int:
    """红色可以都加一 蓝色可以都减一 最小化数组的最大值和最小值的差"""
    red = []  # red 区间可以向右移动
    blue = []  # blue 区间可以向左移动
    for i in range(n):
        if color[i] == "R":
            red.append(nums[i])
        else:
            blue.append(nums[i])

    red.sort()  # [1,2,4]
    blue.sort()  # [3,5]

    if not red:
        return blue[-1] - blue[0]
    if not blue:
        return red[-1] - red[0]

    if red[-1] <= blue[-1]:  # 可以具有包含关系
        blueDiff = blue[-1] - blue[0]
        redDiff = red[-1] - red[0]
        return max(blueDiff, redDiff)

    return red[-1] - min(blue[0], red[0])


t = int(input())
for _ in range(t):
    n = int(input())
    nums = list(map(int, input().split()))
    color = input()
    print(solve1(nums, color))

##############################
# !考虑度数
# n*n矩阵里有(n-2)*(n-2)个点的度数为4
# 剩下的4个点的度数为2
# 4*n-8个点的度数为3
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

n = int(input())


def solve(n: int) -> int:
    res = 4 * (1 + n * n) * n * n // 2
    last = 4 * n - 4
    res -= last * (last + 1) // 2
    res -= 10
    return res % MOD
