from typing import List, Tuple


Circle = Tuple[int, int, int]


# n<=1e4 O(n^2)肯定TLE
#  r <= 10  半径很小

# 注意到半径很小，我们可以枚举每个玩具，并暴力枚举该玩具周围是否有可以套住该玩具的圈。
# !O(r^2*n)
class Solution:
    def circleGame(self, toys: List[List[int]], circles: List[List[int]], r: int) -> int:
        # def isContained(circle1: Circle, circle2: Circle) -> bool:
        #     """若circle1包含于circle2,返回true"""

        #     (x1, y1, r1), (x2, y2, r2) = circle1, circle2
        #     if r1 > r2:
        #         return False
        #     centerDist = ((x1 - x2) ** 2 + (y1 - y2) ** 2) ** 0.5
        #     return centerDist + r1 <= r2

        # res = 0
        # for x1, y1, r1 in toys:
        #     cirCle1 = (x1, y1, r1)
        #     for x2, y2 in circles:
        #         circle2 = (x2, y2, r)
        #         if isContained(cirCle1, circle2):
        #             res += 1
        #             break
        # return res

        def check(tx: int, ty: int, tr: int) -> bool:
            """枚举该玩具周围是否有可以套住该玩具的圈"""
            # 由于圆形区难以枚举，实际枚举的是该圆形区外切的正方形区，要判断一下是否真的是圆心距小于半径差
            for cx in range(tx - r, tx + r + 1):
                for cy in range(ty - r, ty + r + 1):
                    if (cx, cy) not in okSet:
                        continue
                    centerDist = ((cx - tx) ** 2 + (cy - ty) ** 2) ** 0.5
                    if centerDist + tr <= r:
                        return True
            return False

        okSet = set(map(tuple, circles))

        res = 0
        for tx, ty, tr in toys:
            if tr > r:
                continue
            res += int(check(tx, ty, tr))

        return res
