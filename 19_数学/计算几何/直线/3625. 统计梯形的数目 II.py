# 3625. 统计梯形的数目 II
# https://leetcode.cn/problems/count-number-of-trapezoids-ii/description/
# 给你一个二维整数数组 points，其中 points[i] = [xi, yi] 表示第 i 个点在笛卡尔平面上的坐标。
#
# 返回可以从 points 中任意选择四个不同点组成的梯形的数量。
#
# 梯形 是一种凸四边形，具有 至少一对 平行边。两条直线平行当且仅当它们的斜率相同。
#
# 4 <= points.length <= 500
# –1000 <= xi, yi <= 1000
# 所有点两两不同。
#
# !统计直线 + 去掉重复统计的平行四边形
# https://leetcode.cn/problems/count-number-of-trapezoids-ii/solutions/3728529/tong-ji-zhi-xian-qu-diao-zhong-fu-tong-j-a3f9/

from collections import defaultdict
from typing import Counter, List

INF = int(1e18)


class Solution:
    def countTrapezoids(self, points: List[List[int]]) -> int:
        kToB = defaultdict(list)  # 斜率 -> [截距]
        midPointToK = defaultdict(list)  # 中点 -> [斜率]

        for i, (x1, y1) in enumerate(points):
            for x2, y2 in points[:i]:
                dx, dy = x1 - x2, y1 - y2
                k = dy / dx if dx else INF
                b = (y1 * dx - x1 * dy) / dx if dx else x1
                kToB[k].append(b)
                midPointToK[(x1 + x2, y1 + y2)].append(k)

        res = 0
        for bs in kToB.values():
            if len(bs) < 2:
                continue
            pre = 0
            for v in Counter(bs).values():  # !选择一对平行边的方案数
                res += pre * v
                pre += v

        for ks in midPointToK.values():
            if len(ks) < 2:
                continue
            pre = 0
            for v in Counter(ks).values():
                res -= pre * v  # !平行四边形会统计两次，减去多统计的一次
                pre += v

        return res
