# 3609. 到达目标点的最小移动次数(逆向思维 + 分类讨论)
# https://leetcode.cn/problems/minimum-moves-to-reach-target-in-grid/description/
#
# 给你四个整数 sx、sy、tx 和 ty，表示在一个无限大的二维网格上的两个点 (sx, sy) 和 (tx, ty)。
#
# 你的起点是 (sx, sy)。
#
# 在任何位置 (x, y)，定义 m = max(x, y)。你可以执行以下两种操作之一：
#
# 移动到 (x + m, y)，或者
# 移动到 (x, y + m)。
# 返回到达 (tx, ty) 所需的 最小 移动次数。如果无法到达目标点，则返回 -1。
#
# 2543. 判断一个点是否可以到达
# https://leetcode.cn/problems/check-if-point-is-reachable/description/
# 780. 到达终点
# https://leetcode.cn/problems/reaching-points/description/


class Solution:
    def minMoves(self, sx: int, sy: int, tx: int, ty: int) -> int:
        res = 0
        while tx != sx or ty != sy:
            if tx < sx or ty < sy:
                return -1

            if tx == ty:
                if sy > 0:
                    tx = 0
                else:
                    ty = 0
                res += 1
                continue

            # 保证 x > y
            if tx < ty:
                tx, ty = ty, tx
                sx, sy = sy, sx
            if tx > ty * 2:
                if tx % 2 > 0:
                    return -1
                tx //= 2
            else:
                tx -= ty
            res += 1
        return res
