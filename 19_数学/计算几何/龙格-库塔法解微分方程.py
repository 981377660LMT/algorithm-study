# 龙格-库塔法解微分方程 (Runge Kutta)
# 龙格-库塔法是用于非线性常微分方程的解的重要的一类隐式或显式迭代法
# https://kopricky.github.io/code/Misc/runge_kutta.html


from math import cos
from typing import Callable, List, Tuple


class RungeKutta:
    __slots__ = ("_f",)

    def __init__(self, f: Callable[[float, float], float]) -> None:
        """给定一次微分方程 dy/dt=f(t,y), 求解y(t)的值."""
        self._f = f

    def cal(self, t0: float, y0: float, stepSize: float, n: int) -> List[Tuple[float, float]]:
        """runge_kutta法求解微分方程.
        初始值为(t0,y0),步长为stepSize,求解n个点处的值.

        O(n*O(计算一次f))
        """
        res = [(t0, y0)] * (n + 1)
        for i in range(1, n + 1):
            cur = self._proceed(res[i - 1][0], res[i - 1][1], stepSize)
            res[i] = cur
        return res

    def _proceed(self, t: float, y: float, h: float) -> Tuple[float, float]:
        k1 = self._f(t, y)
        k2 = self._f(t + h / 2, y + h / 2 * k1)
        k3 = self._f(t + h / 2, y + h / 2 * k2)
        k4 = self._f(t + h, y + h * k3)
        return t + h, y + h * (k1 + 2 * (k2 + k3) + k4) / 6


if __name__ == "__main__":
    # dy/dt=f(t,y)=y*cos(t)
    def f(t: float, y: float):
        return y * cos(t)

    R = RungeKutta(f)
    print(R.cal(t0=0, y0=1, stepSize=0.5, n=20))
