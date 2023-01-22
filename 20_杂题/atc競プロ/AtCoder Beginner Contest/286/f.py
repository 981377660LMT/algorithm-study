import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# この問題は インタラクティブな問題（あなたが作成したプログラムとジャッジプログラムが標準入出力を介して対話を行う形式の問題）です。

# あなたとジャッジは下記の手順を行います。 手順はフェイズ
# 1 とフェイズ
# 2 からなり、まずフェイズ
# 1 を行った直後、続けてフェイズ
# 2 を行います。

# （フェイズ
# 1 ）

# ジャッジが
# 1 以上
# 10
# 9
#   以下の整数
# N を決める。この整数は隠されている。
# あなたは
# 1 以上
# 110 以下の整数
# M を出力する。
# さらにあなたは、すべての
# i=1,2,…,M について
# 1≤A
# i
# ​
#  ≤M を満たす、長さ
# M の整数列
# A=(A
# 1
# ​
#  ,A
# 2
# ​
#  ,…,A
# M
# ​
#  ) を出力する。
# （フェイズ
# 2 ）

# ジャッジから、長さ
# M の整数列
# B=(B
# 1
# ​
#  ,B
# 2
# ​
#  ,…,B
# M
# ​
#  ) が与えられる。ここで、
# B
# i
# ​
#  =f
# N
#  (i) である。
# f(i) は
# 1 以上
# M 以下の整数
# i に対し
# f(i)=A
# i
# ​
#   で定められ、
# f
# N
#  (i) は
# i を
# f(i) で置き換える操作を
# N 回行った際に得られる整数である。
# あなたは、
# B の情報から、ジャッジが決めた整数
# N を特定し、
# N を出力する。
# 上記の手順を行った後、直ちにプログラムを終了することで正解となります。


# !多个置换环 lcm要大于1e9
# 环的大小由质数构成
# 2 3 5 7 11 13 17 19
# print(2 * 3 * 5 * 7 * 11 * 13 * 17 * 19 * 23 * 5, flush=True)
# print(2 + 3 + 5 + 7 + 11 + 13 + 17 + 19 + 23 * 5, flush=True)

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


if __name__ == "__main__":
    group = []
    start = 1
    for size in [2, 3, 5, 7, 11, 13, 17, 19, 23]:
        cycle = list(range(start, start + size))
        cycle = cycle[1:] + [cycle[0]]
        group.extend(cycle)
        start += size

    # # 第一个5：[7,8,9,10,6]
    # # 第二个5：[11,12,13,14,15]
    # lastGroup = [103, 104, 105, 101, 102]
    # group.extend(lastGroup)

    print(len(group), flush=True)
    print(*group, flush=True)

    f = lambda x: group[x - 1]
    nums = list(map(int, input().split()))  # 長さ M の整数列 B=(B 1 ​ ,B 2 ​ ,…,B M ​ ) が与えられる。

    # n是多少
    remains = []
    start = 0
    for size in [2, 3, 5, 7, 11, 13, 17, 19, 23]:
        target = nums[start : start + size]
        curTime = 0
        while target != group[start : start + size]:
            target = list(map(f, target))
            curTime += 1
        remains.append(curTime)
        start += size

    res = crt(remains, [2, 3, 5, 7, 11, 13, 17, 19, 23])
    print(res + 1, flush=True)  # B から N=2 であると特定しました
