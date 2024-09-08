# No.737 PopCount (bitCountSum)
# https://yukicoder.me/problems/no/737
# 令f(x)=x*x.bit_count()，求f(1)+f(2)+...+f(N)的值模MOD。
# n<=18


import sys
from typing import Tuple


input = lambda: sys.stdin.readline().rstrip("\r\n")
mod = int(1e9 + 7)


def bit_count_sum(higher: int) -> Tuple[int, int]:
    if higher == 0:
        return 0, 0
    if higher % 2:
        a, b = bit_count_sum(higher - 1)
        p = (higher - 1).bit_count()
        a += (higher - 1) * p
        b += p
        return a % mod, b % mod
    a, b = bit_count_sum(higher // 2)
    return (a * 2 + a * 2 + b + (higher // 2) ** 2) % mod, (b * 2 + higher // 2) % mod


if __name__ == "__main__":
    n = int(input())
    res, _ = bit_count_sum(n + 1)
    print(res)
