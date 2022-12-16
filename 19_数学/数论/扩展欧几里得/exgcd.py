from typing import Optional, Tuple


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


assert exgcd(2, 3) == (1, -1, 1)
assert modInv(2, 998244353) == (998244353 + 1) // 2
