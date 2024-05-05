from math import gcd
import sys
from typing import Tuple

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 1 未満の正実数
# r と正整数
# N が与えられます。

# 0≤p≤q≤N かつ
# gcd(p,q)=1 を満たす整数の組
# (p,q) のうち、
# ∣
# ∣
# ∣
# ∣
# ∣
# ​
#  r−
# q
# p
# ​

# ∣
# ∣
# ∣
# ∣
# ∣
# ​
#   の値を最小にするものを求めてください。

# ただし、そのような
# (p,q) が複数存在する場合、
# q
# p
# ​
#   の値が最も小さいものを出力してください。

from typing import Callable, Tuple
from fractions import Fraction


# https://atcoder.jp/contests/abc333/editorial/8000
def sternBrocotTreeSearch(
    n: int, predicate: Callable[[int, int], bool]
) -> Tuple[int, int, int, int]:
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
    r = input()
    n = int(input())
    suffix = r.split(".")[1]
    tmp = len(str(suffix))
    target = Fraction(int(suffix), 10**tmp)
    num, deno = target.numerator, target.denominator
    a, b, c, d = sternBrocotTreeSearch(n, lambda x, y: x * deno < y * num)
    diff1 = abs(Fraction(a, b) - target)
    diff2 = abs(Fraction(c, d) - target)
    if diff1 <= diff2:
        print(f"{a} {b}")
    else:
        print(f"{c} {d}")
