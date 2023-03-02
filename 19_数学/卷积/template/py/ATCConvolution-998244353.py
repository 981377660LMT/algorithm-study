"""
# AtCoder Libary v1.4 を python に移植したもの
# https://github.com/atcoder/ac-library/blob/master/atcoder/convolution.hpp

注意此模板只能用于pypy3 
python3.8会超时,需要numpy的fft
"""

from typing import List


MOD = 998244353
IMAG = 911660635
IIMAG = 86583718
rate2 = (
    0,
    911660635,
    509520358,
    369330050,
    332049552,
    983190778,
    123842337,
    238493703,
    975955924,
    603855026,
    856644456,
    131300601,
    842657263,
    730768835,
    942482514,
    806263778,
    151565301,
    510815449,
    503497456,
    743006876,
    741047443,
    56250497,
    867605899,
    0,
)
irate2 = (
    0,
    86583718,
    372528824,
    373294451,
    645684063,
    112220581,
    692852209,
    155456985,
    797128860,
    90816748,
    860285882,
    927414960,
    354738543,
    109331171,
    293255632,
    535113200,
    308540755,
    121186627,
    608385704,
    438932459,
    359477183,
    824071951,
    103369235,
    0,
)
rate3 = (
    0,
    372528824,
    337190230,
    454590761,
    816400692,
    578227951,
    180142363,
    83780245,
    6597683,
    70046822,
    623238099,
    183021267,
    402682409,
    631680428,
    344509872,
    689220186,
    365017329,
    774342554,
    729444058,
    102986190,
    128751033,
    395565204,
    0,
)
irate3 = (
    0,
    509520358,
    929031873,
    170256584,
    839780419,
    282974284,
    395914482,
    444904435,
    72135471,
    638914820,
    66769500,
    771127074,
    985925487,
    262319669,
    262341272,
    625870173,
    768022760,
    859816005,
    914661783,
    430819711,
    272774365,
    530924681,
    0,
)


def _butterfly(a):
    n = len(a)
    h = (n - 1).bit_length()
    le = 0
    while le < h:
        if h - le == 1:
            p = 1 << (h - le - 1)
            rot = 1
            for s in range(1 << le):
                offset = s << (h - le)
                for i in range(p):
                    l = a[i + offset]
                    r = a[i + offset + p] * rot
                    a[i + offset] = (l + r) % MOD
                    a[i + offset + p] = (l - r) % MOD
                rot *= rate2[(~s & -~s).bit_length()]
                rot %= MOD
            le += 1
        else:
            p = 1 << (h - le - 2)
            rot = 1
            for s in range(1 << le):
                rot2 = rot * rot % MOD
                rot3 = rot2 * rot % MOD
                offset = s << (h - le)
                for i in range(p):
                    a0 = a[i + offset]
                    a1 = a[i + offset + p] * rot
                    a2 = a[i + offset + p * 2] * rot2
                    a3 = a[i + offset + p * 3] * rot3
                    a1na3imag = (a1 - a3) % MOD * IMAG
                    a[i + offset] = (a0 + a2 + a1 + a3) % MOD
                    a[i + offset + p] = (a0 + a2 - a1 - a3) % MOD
                    a[i + offset + p * 2] = (a0 - a2 + a1na3imag) % MOD
                    a[i + offset + p * 3] = (a0 - a2 - a1na3imag) % MOD
                rot *= rate3[(~s & -~s).bit_length()]
                rot %= MOD
            le += 2


def _butterfly_inv(a):
    n = len(a)
    h = (n - 1).bit_length()
    le = h
    while le:
        if le == 1:
            p = 1 << (h - le)
            irot = 1
            for s in range(1 << (le - 1)):
                offset = s << (h - le + 1)
                for i in range(p):
                    l = a[i + offset]
                    r = a[i + offset + p]
                    a[i + offset] = (l + r) % MOD
                    a[i + offset + p] = (l - r) * irot % MOD
                irot *= irate2[(~s & -~s).bit_length()]
                irot %= MOD
            le -= 1
        else:
            p = 1 << (h - le)
            irot = 1
            for s in range(1 << (le - 2)):
                irot2 = irot * irot % MOD
                irot3 = irot2 * irot % MOD
                offset = s << (h - le + 2)
                for i in range(p):
                    a0 = a[i + offset]
                    a1 = a[i + offset + p]
                    a2 = a[i + offset + p * 2]
                    a3 = a[i + offset + p * 3]
                    a2na3iimag = (a2 - a3) * IIMAG % MOD
                    a[i + offset] = (a0 + a1 + a2 + a3) % MOD
                    a[i + offset + p] = (a0 - a1 + a2na3iimag) * irot % MOD
                    a[i + offset + p * 2] = (a0 + a1 - a2 - a3) * irot2 % MOD
                    a[i + offset + p * 3] = (a0 - a1 - a2na3iimag) * irot3 % MOD
                irot *= irate3[(~s & -~s).bit_length()]
                irot %= MOD
            le -= 2


def convolution(s: List[int], t: List[int]) -> List[int]:
    n = len(s)
    m = len(t)
    if min(n, m) <= 60:
        a = [0] * (n + m - 1)
        for i in range(n):
            if i % 8 == 0:
                for j in range(m):
                    a[i + j] += s[i] * t[j]
                    a[i + j] %= MOD
            else:
                for j in range(m):
                    a[i + j] += s[i] * t[j]
        return [x % MOD for x in a]
    a = s.copy()
    b = t.copy()
    z = 1 << (n + m - 2).bit_length()
    a += [0] * (z - n)
    b += [0] * (z - m)
    _butterfly(a)
    _butterfly(b)
    for i in range(z):
        a[i] *= b[i]
        a[i] %= MOD
    _butterfly_inv(a)
    a = a[: n + m - 1]
    iz = pow(z, MOD - 2, MOD)
    return [v * iz % MOD for v in a]


def multiConvolution(arrs: List[List[int]]) -> List[int]:
    if not arrs:
        return []
    if len(arrs) == 1:
        return arrs[0]
    if len(arrs) == 2:
        return convolution(arrs[0], arrs[1])
    m = len(arrs) >> 1
    return convolution(multiConvolution(arrs[:m]), multiConvolution(arrs[m:]))


if __name__ == "__main__":
    MOD = 998244353
    n, m = map(int, input().split())
    a = list(map(int, input().split()))
    b = list(map(int, input().split()))
    print(*convolution(a, b))
