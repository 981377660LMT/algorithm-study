# clamp复合函数
# kind=1 -> f(x) = x + a
# kind=2 -> f(x) = max(x, a)
# kind=3 -> f(x) = min(x, a)
# 对i=1,2,...q 求出复合函数fn(x)=fn(...f2(f1(xi)))的值

# !最后的复合函数形为 F(x) = min(c,max(b,x+a))
from typing import List, Tuple

INF = int(1e18)


Func = Tuple[int, int, int]


def composition(f: "Func", g: "Func") -> "Func":
    a1, b1, c1 = f
    a2, b2, c2 = g
    return a1 + a2, max(b2, b1 + a2), min(c2, max(b2, c1 + a2))


def cal(f: "Func", x: int) -> int:
    a, b, c = f
    return min(c, max(b, x + a))


def clampComposition(funcs: List[Tuple[int, int]], queries: List[int]) -> List[int]:
    f = (0, -INF, INF)  # f(x) = x
    for a, kind in funcs:
        if kind == 1:
            f = composition(f, (a, -INF, INF))
        elif kind == 2:
            f = composition(f, (0, a, INF))
        else:
            f = composition(f, (0, -INF, a))

    return [cal(f, x) for x in queries]


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    funcs = [tuple(map(int, input().split())) for _ in range(n)]  # (a, kind)
    q = int(input())
    queries = list(map(int, input().split()))

    res = clampComposition(funcs, queries)
    print(*res, sep="\n")
