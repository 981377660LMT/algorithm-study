# LCP 74. 最强祝福力场-离散化+二维差分
# https://leetcode.cn/problems/xepqZ5/
# forceField[i] = [x,y,side] 表示第 i 片力场将覆盖以坐标 (x,y) 为中心，边长为 side 的正方形区域。
# !若任意一点的 力场强度 等于覆盖该点的力场数量，请求出在这片地带中 力场强度 最强处的 力场强度。


from typing import List
from 二维差分模板 import DiffMatrix


class Solution:
    def fieldOfGreatestBlessing(self, forceField: List[List[int]]) -> int:
        ...
