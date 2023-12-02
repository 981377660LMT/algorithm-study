from typing import Callable, Tuple


def sternBrocotTreeSearch(
    n: int, predicate: Callable[[int, int], bool]
) -> Tuple[int, int, int, int]:
    """
    返回分子和分母都不超过n的最简分数p/q中,
    满足条件predicate(p/q)的最大值a/b,以及不满足条件predicate(p/q)的最小值c/d。
    predicate(p/q)是单调的。
    时间复杂度为O(f(n)logn),其中f(n)是计算predicate(p/q)的时间复杂度。
    """
    a, b, c, d = 0, 1, 1, 0
    while True:
        num = a + c
        den = b + d
        if num > n or den > n:
            break
        if predicate(num, den):
            k = 2
            while True:
                num = a + k * c
                if num > n:
                    break
                den = b + k * d
                if den > n:
                    break
                if not predicate(num, den):
                    break
                k *= 2
            k //= 2
            a += c * k
            b += d * k
        else:
            k = 2
            while True:
                num = a * k + c
                if num > n:
                    break
                den = b * k + d
                if den > n:
                    break
                if predicate(num, den):
                    break
                k *= 2
            k //= 2
            c += a * k
            d += b * k
    return a, b, c, d


if __name__ == "__main__":
    # https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=1208
    while True:
        a, b = map(int, input().split())
        if a == 0 and b == 0:
            break
        res = sternBrocotTreeSearch(b, lambda x, y: x * x < y * y * a)
        print(f"{res[2]}/{res[3]} {res[0]}/{res[1]}")
