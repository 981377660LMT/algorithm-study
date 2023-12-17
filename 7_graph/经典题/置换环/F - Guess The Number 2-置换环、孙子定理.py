"""
置换环+中国剩余定理(孙子定理)

你需要猜到对方的n的大小,你仅可做一件事：

- 指定m,输出一个长度为 m的数组A。其中 1≤m≤110,1≤Ai≤m
- 对方根据数组 A,输出一个数组 B
- 你根据该数组 B,得出 n的值。

生成数组 B的方式为:
假想一个有m个点的图,其中点i连一条有向边到点ai。
Bi的值为:从i号点出发,走 n条边到达的点的编号。

n≤1e9

解:
# !多个置换环的大小 lcm要大于1e9
# 环的大小需要由两两互素的数构成
# 4,9,5,7,11,13,17,19,23,其和为108<110,其乘积为1338557220>1e9
通过该数组(每个元素是一个环的大小)可以构造出对应的数组a。
然后根据数组b计算得到每个余数 ri,然后由中国剩余定理解出 n。
"""


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


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


# def modInv(a: int, mod: int) -> int:
#     return pow(a, mod - 2, mod)


def modInv(a: int, mod: int) -> Optional[int]:
    """
    扩展gcd求a在mod下的逆元
    即求出逆元 `inv` 满足 `a*inv ≡ 1 (mod m)`
    """
    gcd_, x, _ = exgcd(a, mod)
    if gcd_ != 1:
        return None
    return x % mod


if __name__ == "__main__":
    MODS = [4, 9, 5, 7, 11, 13, 17, 19, 23]
    A = []
    start = 1
    starts = []  # 每个置换环的起点
    for size in MODS:
        cycle = list(range(start, start + size))
        cycle = cycle[1:] + [cycle[0]]
        A.extend(cycle)
        starts.append(start)
        start += size

    print(len(A), flush=True)
    print(*A, flush=True)

    B = list(map(int, input().split()))  # 長さ M の整数列 B=(B 1 ​ ,B 2 ​ ,…,B M ​ ) が与えられる。

    # 求n
    remains = [B[start - 1] - start for start in starts]
    res = crt(remains, MODS)
    print(res, flush=True)
