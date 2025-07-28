# 士兵可以进行移动，每次移动，一名士兵可以向上，向下，向左或向右移动一个单位
# （因此，他的 x 或 y 坐标也将加 1 或减 1）。
# 现在希望通过移动士兵，使得所有士兵彼此相邻的处于同一条水平线内，
# 即所有士兵的 y 坐标相同并且 x 坐标相邻。
# 请你计算满足要求的情况下，所有士兵的总移动次数最少是多少。
# 需注意，两个或多个士兵不能占据同一个位置。
#
#
# 1. 上下移动与左右移动可以分开进行
# 2. 找出y坐标的中位数，用最小的代价先将所有点移到同一列上
# 3. 将这些点彼此相邻
#
# Q4. 棋盘整理
# https://leetcode.cn/contest/2025_pudong_ai/problems/1Hxnb6/description/
# 考虑将 x 值变得相等的最少距离。
# 这个问题相当于：给很多数，每次可以将一个数增加或减少 1，问将所有数变得相等的最少次数。结论是都变成它们的中位数。
# 考虑将 y 值变得连续的最少距离。这个问题相当于：给很多数，每次可以将一个数增加或减少 1，问将所有数变得连续的最少次数。
# 结论是先把所有数从小到大排序，然后第 i 个数减去 i，就变成了前一个问题。

from typing import List


def solve1(nums: List[int]) -> int:
    nums = sorted(nums)
    n = len(nums)
    mid = nums[n // 2]
    return sum(abs(x - mid) for x in nums)


def solve2(nums: List[int]) -> int:
    nums = sorted(nums)
    for i in range(len(nums)):
        nums[i] -= i
    return solve1(nums)


def solve(xs: List[int], ys: List[int]) -> int:
    return solve1(xs) + solve2(ys)


class Solution:
    def organizeChessboard(self, pieces: List[List[int]]) -> int:
        xs, ys = [x for x, _ in pieces], [y for _, y in pieces]
        res1 = solve(xs, ys)
        res2 = solve(ys, xs)
        return min(res1, res2) % int(1e9 + 7)
