"""
求RGB串的个数
R G B 的个数分别为R G B 个 (R,G,B <=1e6)
RG子串的个数为K个


https://atcoder.jp/contests/abc266/editorial/4721
将RG视为K 则等价于
K个K (R-K)个R (G-K)个G B个B 组成的没有RG子串的字符串个数
# !1.先安排 RG、R、B 随意排列
# !2.然后剩下的G安排到槽里(注意G不能放在R后的槽里)
"""


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


MOD = 998244353
fac = [1]
ifac = [1]
for i in range(1, int(2e6) + 10):  # ! G + B 最多2e6 (注意3e6会TLE)
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


def A(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return (fac[n] * ifac[n - k]) % MOD


_r, _g, _b, _k = map(int, input().split())

K, R, G, B = _k, _r - _k, _g - _k, _b
slots = K + B + 1


res1 = fac[K + B + R] * ifac[K] % MOD * ifac[B] % MOD * ifac[R] % MOD
res2 = C(slots + G - 1, slots - 1)
print(res1 * res2 % MOD)
