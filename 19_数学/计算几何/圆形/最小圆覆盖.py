# 你需要用最少的原材料给花园安装一个 圆形 的栅栏，
# 使花园中所有的树都在被 围在栅栏内部（在栅栏边界上的树也算在内）。
# https://leetcode.cn/problems/erect-the-fence-ii/
# !最小圆覆盖 Welzl 算法(随机增量法)
# n<=3000


import random
from typing import List, Tuple, Union


EPS = 1e-6


def calCircle2(
    x1: int, y1: int, x2: int, y2: int, x3: int, y3: int
) -> Union[Tuple[None, None, None], Tuple[float, float, float]]:
    """三点圆公式,求圆的圆心坐标(x,y)和半径r"""
    a, b, c, d = x1 - x2, y1 - y2, x1 - x3, y1 - y3
    a1 = (x1 * x1 - x2 * x2 + y1 * y1 - y2 * y2) / 2
    a2 = (x1 * x1 - x3 * x3 + y1 * y1 - y3 * y3) / 2
    theta = b * c - a * d
    if theta == 0:
        return None, None, None
    x0 = (b * a2 - d * a1) / theta
    y0 = (c * a1 - a * a2) / theta
    return x0, y0, ((x1 - x0) * (x1 - x0) + (y1 - y0) * (y1 - y0)) ** 0.5


def isCover(x1: int, y1: int, rx: float, ry: float, r: float) -> bool:
    dist = ((x1 - rx) * (x1 - rx) + (y1 - ry) * (y1 - ry)) ** 0.5
    return dist <= r or abs(dist - r) < EPS


class Solution:
    def outerTrees(self, trees: List[List[int]]) -> List[float]:
        """请用一个长度为 3 的数组 [x,y,r] 来返回最小圆的圆心坐标和半径

        如果答案与正确答案的误差不超过 1e-5,则该答案将被视为正确答案通过
        """
        random.shuffle(trees)

        n = len(trees)
        x0, y0 = trees[0]
        r = 0
        for i in range(1, n):
            x1, y1 = trees[i]
            if isCover(x1, y1, x0, y0, r):
                continue
            x0, y0, r = x1, y1, 0
            for j in range(i):
                x2, y2 = trees[j]
                if isCover(x2, y2, x0, y0, r):
                    continue
                x0, y0, r = (
                    (x1 + x2) / 2,
                    (y1 + y2) / 2,
                    (((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2)) ** 0.5) / 2,
                )
                for k in range(j):
                    x3, y3 = trees[k]
                    if isCover(x3, y3, x0, y0, r):
                        continue
                    candX, candY, candR = calCircle2(x1, y1, x2, y2, x3, y3)
                    if candX is not None and candY is not None and candR is not None:
                        x0, y0, r = candX, candY, candR

        return [x0, y0, r]


if __name__ == "__main__":
    print(Solution().outerTrees([[1, 1], [2, 2], [2, 0], [2, 4], [3, 3], [4, 2]]))
    print(
        Solution().outerTrees(
            [
                [55, 7],
                [36, 30],
                [1, 64],
                [83, 97],
                [8, 90],
                [16, 7],
                [18, 23],
                [98, 77],
                [57, 33],
                [98, 54],
                [73, 7],
            ]
        )
    )
