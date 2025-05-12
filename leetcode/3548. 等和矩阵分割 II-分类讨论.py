# 3548. 等和矩阵分割 II
# https://leetcode.cn/problems/equal-sum-grid-partition-ii/description/
# 给你一个由正整数组成的 m x n 矩阵 grid。你的任务是判断是否可以通过 一条水平或一条垂直分割线 将矩阵分割成两部分，使得：
#
# 分割后形成的每个部分都是 非空 的。
# 两个部分中所有元素的和 相等 ，或者总共 最多移除一个单元格 （从其中一个部分中）的情况下可以使它们相等。
# 如果移除某个单元格，剩余部分必须保持 连通 。
# 如果存在这样的分割，返回 true；否则，返回 false。
#
# 特判 n=1 的情况，此时只能删除第一个数或者分割线上那个数。
# 如果分割线在第一行与第二行之间，那么不能删除第一行的中间元素（首尾元素可以删），因为删除后会导致第一部分不连通。
# 其余情况，可以随便删。


from typing import List


class Solution:
    def canPartitionGrid(self, grid: List[List[int]]) -> bool:
        total = sum(map(sum, grid))

        def check(grid: List[List[int]]) -> bool:
            COL = len(grid[0])

            def f() -> bool:
                visited = {0}
                curSum = 0
                for r, row in enumerate(grid):
                    for c, v in enumerate(row):
                        curSum += v
                        if r > 0 or c == 0 or c == COL - 1:
                            visited.add(v)
                    if COL == 1:
                        if (
                            curSum * 2 == total
                            or curSum * 2 - total == grid[0][0]
                            or curSum * 2 - total == row[0]
                        ):
                            return True
                        continue
                    if curSum * 2 - total in visited:
                        return True
                    if r == 0:
                        visited.update(row)
                return False

            if f():  # 最多删除上半中的一个数
                return True
            grid.reverse()
            return f()  # 最多删除下半中的一个数

        return check(grid) or check(list(zip(*grid))[::-1])  # type: ignore
