# 3189. 得到一个和平棋盘的最少步骤
# https://leetcode.cn/problems/minimum-moves-to-get-a-peaceful-board/description/
#
# 给定一个长度为 n 的二维数组 rooks，其中 rooks[i] = [xi, yi] 表示 n x n 棋盘上一个车的位置。你的任务是每次在垂直或水平方向上移动 1 格 车（到一个相邻的格子）使得棋盘变得 和平。
# 如果每行每列都 只有 一个车，那么这块棋盘就是和平的。
# 返回获得和平棋盘所需的 最少 步数。
# 注意 任何时刻 两个车都不能在同一个格子。
#
# 如果存在y挡住了x前往目标点x'，那么先将y移动到目标的y'，
# 如果y目标点y'被x挡住了，那么交换xy的目标点不会增加步数，
# !因此不需要考虑阻挡的情况，只需要考虑最终结果。
# 二维互不干扰互相独立，可以当成两次一维情况来考虑，即先按横坐标排序，再按纵坐标排序，得到对应结果

from typing import List


class Solution:
    def minMoves(self, rooks: List[List[int]]) -> int:
        def solve(nums: List[int]) -> int:
            n = len(nums)
            nums = sorted(nums)
            return sum(abs(nums[i] - i) for i in range(n))

        resX = solve([x for x, _ in rooks])
        resY = solve([y for _, y in rooks])
        return resX + resY
