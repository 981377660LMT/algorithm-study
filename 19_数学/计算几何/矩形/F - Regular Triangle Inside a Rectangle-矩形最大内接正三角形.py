# F - Regular Triangle Inside a Rectangle-矩形最大内接正三角形边长
from math import sqrt


def solve1(a: int, b: int) -> float:
    SQRT3 = 3**0.5
    if a < b:
        a, b = b, a
    if a >= 2 * b / SQRT3:
        s = b * b / SQRT3
        return sqrt(s * 4 / SQRT3)
    s = SQRT3 * a * a - 3 * a * b + SQRT3 * b * b
    return sqrt(s * 4 / SQRT3)


if __name__ == "__main__":
    a, b = map(int, input().split())
    print(solve1(a, b))
