# 解一元二次方程/解方程
# !如果需要避免浮点数误差,需要用Decimal运算

from math import sqrt
from typing import List, Tuple


def quadratic_solver(A: int, B: int, C: int) -> Tuple[int, List[float]]:
    """A*x^2 + B*x + C = 0

    返回解的个数和解的列表(-1表示无穷多个解)
    """
    if B < 0:
        A = -A
        B = -B
        C = -C
    if A == 0:
        if B == 0:
            if C == 0:
                return -1, []  # 无穷多个解
            return 0, []  # 无解
        return 1, [-C / B]

    D = B * B - 4 * A * C
    if D < 0:
        return 0, []
    if D == 0:
        return 1, [-B / (2 * A)]
    res1 = (-B - sqrt(D)) / (2 * A)
    res2 = C / A / res1
    if res1 > res2:
        res1, res2 = res2, res1
    return 2, [res1, res2]


if __name__ == "__main__":
    # https://yukicoder.me/problems/no/955%3E
    a, b, c = map(int, input().split())
    n, res = quadratic_solver(a, b, c)
    if n == -1:
        print(-1)
    else:
        print(n)
        for x in res:
            print(x)
