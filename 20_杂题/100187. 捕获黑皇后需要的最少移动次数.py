# 100187. 捕获黑皇后需要的最少移动次数
# https://leetcode.cn/problems/minimum-moves-to-capture-the-queen/description/
# 现有一个下标从 0 开始的 8 x 8 棋盘，上面有 3 枚棋子。

# 给你 6 个整数 a 、b 、c 、d 、e 和 f ，其中：

# (a, b) 表示白色车的位置。
# (c, d) 表示白色象的位置。
# (e, f) 表示黑皇后的位置。
# 假定你只能移动白色棋子，返回捕获黑皇后所需的最少移动次数。

# 请注意：
# 车可以向垂直或水平方向移动任意数量的格子，但不能跳过其他棋子。
# 象可以沿对角线方向移动任意数量的格子，但不能跳过其他棋子。
# 如果车或象能移向皇后所在的格子，则认为它们可以捕获皇后。
# 皇后不能移动。


# 答案为1或2
# 车和皇后在同一行或同一列 => 1 or 2
# 象和皇后在同一对角线 => 1 or 2
# 其余情况 => 2
class Solution:
    def minMovesToCaptureTheQueen(self, a: int, b: int, c: int, d: int, e: int, f: int) -> int:
        if a == e:
            min_, max_ = min(b, f), max(b, f)
            if min_ < d < max_ and c == e:
                return 2
            else:
                return 1
        if b == f:
            min_, max_ = min(a, e), max(a, e)
            if min_ < c < max_ and d == f:
                return 2
            else:
                return 1
        if abs(c - e) == abs(d - f):
            minx, maxx = min(c, e), max(c, e)
            miny, maxy = min(d, f), max(d, f)
            if minx < a < maxx and miny < b < maxy and abs(a - c) == abs(b - d):
                return 2
            else:
                return 1
        return 2
