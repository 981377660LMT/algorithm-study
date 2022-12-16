# !https://github.dev/EndlessCheng/codeforces-go/blob/621c08102b20f4039664ce41972ce65a0daaad33/copypasta/math.go#L1873
# 中国剩余定理(孙子定理)
# Chinese Remainder Theorem, CRT


from math import gcd
from typing import List, Optional, Tuple


def crt(remains: List[int], mods: List[int]) -> Optional[int]:
    """
    `模数两两互素`的线性同余方程组的最小非负整数解 - 中国剩余定理 (CRT)
    x ≡ remains_i (mod mods_i), mods_i 两两互质且 Πmods_i <= 1e18
    """
    modMul = 1
    for m in mods:
        modMul *= m
    res = 0
    for mod, remain in zip(mods, remains):
        other = modMul // mod
        inv = modInv(other, mod)
        if inv is None:
            return None
        res = (res + remain * other * inv) % modMul
    return res


def excrt(A: List[int], remains: List[int], mods: List[int]) -> Optional[Tuple[int, int]]:
    """
    线性同余方程组的最小非负整数解 - 扩展中国剩余定理 (EXCRT)
    A_i * x ≡ remains_i (mod mods_i), Πmods_i <= 1e18

    Returns:
      Optional[Tuple[int, int]]:
        解为 x ≡ b (mod m)
        有解时返回 (b, m),无解时返回None
    """
    modMul = 1
    res = 0
    for i, mod in enumerate(mods):
        a, b = A[i] * modMul, remains[i] - A[i] * res
        d = gcd(a, mod)
        if b % d != 0:
            return None
        t = rationalMod(b // d, a // d, mod // d)
        if t is None:
            return None
        res += modMul * t
        modMul *= mod // d
    return res % modMul, modMul


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


def modInv(a: int, mod: int) -> Optional[int]:
    """
    扩展gcd求a在mod下的逆元
    即求出逆元 `inv` 满足 `a*inv ≡ 1 (mod m)`
    """
    gcd_, x, _ = exgcd(a, mod)
    if gcd_ != 1:
        return None
    return x % mod


def rationalMod(a: int, b: int, mod: int) -> Optional[int]:
    """
    有理数取模(有理数取余)
    求 a/b 模 mod 的值
    """
    inv = modInv(b, mod)
    if inv is None:
        return None
    return a * inv % mod


assert crt([2, 3, 2], [3, 5, 7]) == 23
assert excrt([1, 1, 1], [2, 3, 2], [3, 5, 7]) == (23, 105)
