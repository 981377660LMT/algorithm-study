# 3453. 分割正方形 I
# https://leetcode.cn/problems/separate-squares-i/
# 给你一个二维整数数组 squares ，其中 squares[i] = [xi, yi, li] 表示一个与 x 轴平行的正方形的左下角坐标和正方形的边长。
# 找到一个最小的 y 坐标，它对应一条水平线，该线需要满足它以上正方形的总面积 等于 该线以下正方形的总面积。
# 答案如果与实际答案的误差在 10-5 以内，将视为正确答案。
# 注意：正方形 可能会 重叠。重叠区域应该被 多次计数 。

from math import ceil
from typing import Callable, List


class Solution:
    def separateSquares(self, squares: List[List[int]]) -> float:
        EPS_INV = int(1e5)
        areaSum = sum(l * l for _, _, l in squares)

        def check(mid: float) -> bool:
            sum_ = 0
            for _, y, l in squares:
                if y < mid:
                    sum_ += l * min(mid - y, l)
            return sum_ >= areaSum / 2

        minY = min(y for _, y, _ in squares)
        maxY = max(y + l for _, y, l in squares)
        return bisectLeftFloat(minY, maxY, check, EPS_INV)


def bisectLeftFloat(
    left: float, right: float, check: Callable[[float], bool], epsInv=int(1e9)
) -> float:
    diff = ceil((right - left) * epsInv)
    round = diff.bit_length()
    for _ in range(round):
        mid = (left + right) / 2
        if check(mid):
            right = mid
        else:
            left = mid
    return (left + right) / 2


def bisectRightFloat(
    left: float, right: float, check: Callable[[float], bool], epsInv=int(1e9)
) -> float:
    diff = ceil((right - left) * epsInv)
    round = diff.bit_length()
    for _ in range(round):
        mid = (left + right) / 2
        if check(mid):
            left = mid
        else:
            right = mid
    return (left + right) / 2
