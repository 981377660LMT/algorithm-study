# 3858. 按位或的最小值
# https://leetcode.cn/problems/minimum-bitwise-or-from-grid/description/
# 给你一个大小为 m x n 的二维整数数组 grid。
# 你必须从 grid 的每一行中 选择恰好一个整数。
# 返回一个整数，表示从每行中选出的整数的 按位或（bitwise OR）的 最小可能值。
#
# 试填法
# 要想让答案尽量小，那么答案二进制的高位是 0 比是 1 更好，所以优先判断答案的高位能否是 0，
# 即从高到低依次判断答案的第 d 位能不能是 0

from typing import List


class Solution:
    def minimumOR(self, grid: List[List[int]]) -> int:
        max_ = max(max(row) for row in grid)
        res = 0
        for d in range(max_.bit_length() - 1, -1, -1):
            mask = res | ((1 << d) - 1)
            for row in grid:
                for x in row:
                    if x | mask == mask:
                        break
                else:
                    res |= 1 << d
                    break
        return res
