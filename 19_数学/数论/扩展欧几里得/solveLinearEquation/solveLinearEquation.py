from typing import Tuple


def solveLinearEquation(a: int, b: int, c: int, allowZero=False) -> Tuple[int, int, int, int, int]:
    """
    a*x + b*y = c 的通解为
    x = (c/g)*x0 + (b/g)*k
    y = (c/g)*y0 - (a/g)*k
    其中 g = gcd(a,b) 且需要满足 g|c
    x0 和 y0 是 ax+by=g 的一组特解（即 exgcd(a,b) 的返回值）

    为方便讨论，这里要求输入的 a b c 必须为正整数

    返回: 正整数解的个数（无解时为 -1,无正整数解时为 0)
          x 取最小正整数时的解 x1 y1,此时 y1 是最大正整数解
          y 取最小正整数时的解 x2 y2,此时 x2 是最大正整数解
    """
    g, x0, y0 = exgcd(a, b)

    # 无解
    if c % g != 0:
        return -1, 0, 0, 0, 0

    a //= g
    b //= g
    c //= g
    x0 *= c
    y0 *= c

    x1 = x0 % b
    if allowZero:
        if x1 < 0:
            x1 += b
    else:
        if x1 <= 0:
            x1 += b
    k1 = (x1 - x0) // b
    y1 = y0 - k1 * a

    y2 = y0 % a
    if allowZero:
        if y2 < 0:
            y2 += a
    else:
        if y2 <= 0:
            y2 += a
    k2 = (y0 - y2) // a
    x2 = x0 + k2 * b

    # 无正整数解
    if y1 <= 0:
        return 0, x1, y1, x2, y2

    # k 越大 x 越大
    return k2 - k1 + 1, x1, y1, x2, y2


def exgcd(a: int, b: int) -> Tuple[int, int, int]:
    """求解二元一次不定方程 ax+by=gcd(a,b) 的特解(x,y),特解满足 |x|<=|b|, |y|<=|a|."""
    if b == 0:
        return a, 1, 0
    gcd_, y, x = exgcd(b, a % b)
    y -= a // b * x
    return gcd_, x, y


if __name__ == "__main__":
    # P5656 【模板】二元一次不定方程 (exgcd)
    # https://www.luogu.com.cn/problem/P5656
    def p5656() -> None:
        T = int(input())
        for _ in range(T):
            a, b, c = map(int, input().split())
            n, x1, y1, x2, y2 = solveLinearEquation(a, b, c, False)
            if n == -1:
                print(-1)
            elif n == 0:
                print(x1, y2)
            else:
                print(n, x1, y2, x2, y1)

    # https://www.luogu.com.cn/problem/CF1244C
    # (x,y,z)表示一个为非负整数三元组，满足
    # x+y+z=w，a*x+b*y=c
    # 无解输出 -1，否则输出任意一组正整数解.
    # !让非负解 x+y 尽量小,最简单的做法就是 min(x1+y1, x2+y2).
    def cf1244c() -> None:
        w, c, a, b = map(int, input().split())
        n, x1, y1, x2, y2 = solveLinearEquation(a, b, c, True)
        if n == -1:
            print(-1)
        else:
            res1, res2 = x1, y1
            if res1 < 0 or res2 < 0 or res1 + res2 > x2 + y2:
                res1, res2 = x2, y2
            if res1 + res2 > w or res1 < 0 or res2 < 0:
                print(-1)
            else:
                print(res1, res2, w - res1 - res2)

    cf1244c()
