from dataclasses import dataclass
from math import cos, pi, sin, sqrt
from random import random, uniform
from typing import List


@dataclass
class Solution:
    radius: float
    x_center: float
    y_center: float

    def randPoint(self) -> List[float]:
        # !在单位圆中随机生成一个点，它离圆心的距离小于等于 r 的概率为 F(r) = r^2
        randR = self.radius * sqrt(random())  # 注意这里是根号
        randArc = uniform(0, 2 * pi)
        x = self.x_center + randR * cos(randArc)
        y = self.y_center + randR * sin(randArc)
        return [x, y]

