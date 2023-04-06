# P3187-最小矩形覆盖
# 给定一些点的坐标，要求求能够覆盖所有点的最小面积的矩形，输出所求矩形的面积和四个顶点坐标
# https://www.luogu.com.cn/problem/P3187
# https://www.luogu.com.cn/record/42463169


from typing import List, Tuple

Point = Tuple[float, float]


def minRectangle(points: List[Point]) -> Tuple[int, Tuple[Point, Point, Point, Point]]:
    ...
