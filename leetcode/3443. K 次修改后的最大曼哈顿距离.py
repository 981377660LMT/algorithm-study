# 3443. K 次修改后的最大曼哈顿距离
# https://leetcode.cn/problems/maximum-manhattan-distance-after-k-changes/description/
# 给你一个由字符 'N'、'S'、'E' 和 'W' 组成的字符串 s，其中 s[i] 表示在无限网格中的移动操作：
# 'N'：向北移动 1 个单位。
# 'S'：向南移动 1 个单位。
# 'E'：向东移动 1 个单位。
# 'W'：向西移动 1 个单位。
# 初始时，你位于原点 (0, 0)。你 最多 可以修改 k 个字符为任意四个方向之一。
# 请找出在 按顺序 执行所有移动操作过程中的 任意时刻 ，所能达到的离原点的 最大曼哈顿距离 。
# 曼哈顿距离 定义为两个坐标点 (xi, yi) 和 (xj, yj) 的横向距离绝对值与纵向距离绝对值之和，即 |xi - xj| + |yi - yj|。
#
# !每操作一次，曼哈顿距离都会增大 2，但这不会超过移动的次数 i+1


class Solution:
    def maxDistance(self, s: str, k: int) -> int:
        res, x, y = 0, 0, 0
        for i, c in enumerate(s):
            if c == "N":
                y += 1
            elif c == "S":
                y -= 1
            elif c == "E":
                x += 1
            else:
                x -= 1
            v1 = abs(x) + abs(y) + 2 * k
            v2 = i + 1
            res = max(res, min(v1, v2))
        return res
