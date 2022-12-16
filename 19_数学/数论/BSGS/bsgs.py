"""bsgs与exbsgs

https://dianhsu.com/2022/08/27/template-math/#bsgs
"""


from math import ceil, gcd, sqrt
from typing import Tuple


def bsgs(base: int, target: int, p: int) -> int:
    """Baby-step Giant-step

    在base和p互质的情况下,求解 base^x ≡ target (mod p) 的最小解x,
    若不存在解则返回-1

    时间复杂度: O(sqrt(p)))

    https://dianhsu.com/2022/08/27/template-math/#bsgs
    """
    mp = dict()
    t = ceil(sqrt(p))
    target %= p
    val = 1
    for i in range(t):
        tv = target * val % p
        mp[tv] = i
        val = val * base % p

    base, val = val, 1
    if base == 0:
        return 1 if target == 0 else -1

    for i in range(t + 1):
        tv = mp.get(val, -1)
        if tv != -1 and i * t - tv >= 0:
            return i * t - tv
        val = val * base % p

    return -1


def exgcd(a: int, b: int) -> Tuple[int, int, int]:
    """
    求a, b最大公约数,同时求出裴蜀定理中的一组系数x, y,
    满足 x*a + y*b = gcd(a, b)

    ax + by = gcd_ 返回 `(gcd_, x, y)`
    """
    if b == 0:
        return a, 1, 0
    gcd_, x, y = exgcd(b, a % b)
    return gcd_, y, x - a // b * y


def exbsgs(base: int, target: int, p: int) -> int:
    """Extended Baby-step Giant-step

    求解 base^x ≡ target (mod p) 的最小解x,
    若不存在解则返回-1

    时间复杂度: O(sqrt(p)))

    https://dianhsu.com/2022/08/27/template-math/#exbsgs
    """
    base %= p
    target %= p

    # !平凡解
    if target == 1 or p == 1:
        return 0

    cnt = 0
    d, ad = 1, 1
    while True:
        d = gcd(base, p)
        if d == 1:
            break
        if target % d:
            return -1
        cnt += 1
        target //= d
        p //= d
        ad = ad * (base // d) % p
        if ad == target:
            return cnt

    gcd_, x, _y = exgcd(ad, p)
    inv = x % p
    res = bsgs(base, target * inv % p, p)
    if res != -1:
        res += cnt
    return res


if __name__ == "__main__":
    # https://www.luogu.com.cn/problem/P4195
    # !给定a,p,b，求满足a**x ≡ b (mod p)的最小自然数x。

    while True:
        base, p, target = map(int, input().split())
        if base == target == p == 0:
            break
        res = exbsgs(base, target, p)
        print(res if res != -1 else "No Solution")
