"""bsgs与exbsgs 离散对数(Discrete Logarithm)

https://dianhsu.com/2022/08/27/template-math/#bsgs
"""

from math import ceil, gcd, sqrt
from typing import Tuple


def bsgs(base: int, target: int, mod: int) -> int:
    """Baby-step Giant-step

    在base和mod互质的情况下,求解 base^x ≡ target (mod mod) 的最小解x,
    若不存在解则返回-1

    时间复杂度: O(sqrt(mod)))

    https://dianhsu.com/2022/08/27/template-math/#bsgs
    """
    base %= mod
    target %= mod
    if target == 1 or mod == 1:
        return 0

    mp = dict()
    t = ceil(sqrt(mod))
    val = 1
    for i in range(t):
        tv = target * val % mod
        mp[tv] = i
        val = val * base % mod

    base, val = val, 1
    if base == 0:
        return 1 if target == 0 else -1

    for i in range(t + 1):
        tv = mp.get(val, -1)
        if tv != -1 and i * t - tv >= 0:  # !注意这里取等号表示允许最小解为0
            return i * t - tv
        val = val * base % mod

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
    if target == 1 or p == 1:  # !注意这里允许最小解为0
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

    _, x, _ = exgcd(ad, p)
    inv = x % p
    res = bsgs(base, target * inv % p, p)
    if res != -1:
        res += cnt
    return res


if __name__ == "__main__":
    # https://judge.yosupo.jp/problem/discrete_logarithm_mod
    T = int(input())
    for _ in range(T):
        base, target, mod = map(int, input().split())
        res = exbsgs(base, target, mod)
        print(res)

    # https://www.luogu.com.cn/problem/P4195
    # !给定a,p,b，求满足a**x ≡ b (mod p)的最小自然数x。
    while True:
        base, p, target = map(int, input().split())
        if base == target == p == 0:
            break
        res = exbsgs(base, target, p)
        print(res if res != -1 else "No Solution")

    class Solution:
        def minAllOneMultiple1(self, k: int) -> int:
            """
            求解 10^x = 1 (mod m) 的最小正整数解
            为了避开 x=0, 求解 10^{x-1} = 10^-1 (mod m)
            """
            if k % 2 == 0 or k % 5 == 0:
                return -1
            m = 9 * k
            inv10 = pow(10, -1, m)
            res = bsgs(10, inv10, m)
            return res + 1 if res != -1 else -1

        def minAllOneMultiple2(self, k: int) -> int:
            if k % 2 == 0 or k % 5 == 0:
                return -1
            mod_, len_ = 0, 1
            while True:
                mod_ = (mod_ * 10 + 1) % k
                if mod_ == 0:
                    return len_
                len_ += 1
