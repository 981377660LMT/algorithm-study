# multi Dimesional Selection

# !二维矩阵选数的最小值 残酷群OA
# 从各行选一些数 每行至少选择ceil(COL/2)个数
# 选择的代价为所有数的最大值减最小值
# 设k为所有`选择代价的最小值`
# !求`选择代价的最小值`与`代价最小时选择的个数` 的乘积最大值
# ROW,COL<=1000

# 分析:类似16_滑动窗口/k重复字符子串词频统计/面试题 17.18. 最短超串.py
# !答案肯定是排序后的一段子串
# 时间复杂度O(mnlogmn)

from math import ceil
from typing import List

INF = int(1e18)


def getMaxProduct(n: int, m: int, grid: List[List[int]]) -> int:
    ROW, COL = n, m
    nums = []  # !(value,row)
    for r in range(ROW):
        for c in range(COL):
            nums.append((grid[r][c], r))
    nums.sort(key=lambda x: x[0])

    rowNeed = ceil(COL / 2)
    res, left, n = 0, 0, len(nums)
    rowCounter, okRow, minDiff = [0] * ROW, 0, INF
    for right in range(n):
        rowCounter[nums[right][1]] += 1
        if rowCounter[nums[right][1]] == rowNeed:
            okRow += 1

        while left <= right and okRow == ROW:
            diff = nums[right][0] - nums[left][0]
            if diff == minDiff:
                cand = minDiff * (right - left + 1)
                res = cand if cand > res else res
            elif diff < minDiff:
                minDiff = diff
                res = minDiff * (right - left + 1)
            rowCounter[nums[left][1]] -= 1
            if rowCounter[nums[left][1]] == rowNeed - 1:
                okRow -= 1
            left += 1

    return res


assert getMaxProduct(3, 2, [[1, 2], [3, 4], [8, 9]]) == 24
# 最小值为6 对应[2,4,8]和[2,3,4,8]
