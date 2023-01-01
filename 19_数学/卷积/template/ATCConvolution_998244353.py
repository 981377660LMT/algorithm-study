"""
Reference
https://github.com/atcoder/ac-library/blob/master/atcoder/convolution.hpp
https://github.com/atcoder/ac-library/blob/master/atcoder/internal_math.hpp
https://github.com/atcoder/ac-library/blob/master/document_en/convolution.md
https://github.com/atcoder/ac-library/blob/master/document_ja/convolution.md

注意此模板只能用于pypy3 
python3.8会超时,需要numpy的fft
"""
MOD = 998244353


def primitive_root(m):
    if m == 2:
        return 1
    if m == 167772161:
        return 3
    if m == 469762049:
        return 3
    if m == 754974721:
        return 11
    if m == 998244353:
        return 3
    divs = [0] * 20
    divs[0] = 2
    cnt = 1
    x = (m - 1) // 2
    while x % 2 == 0:
        x //= 2
    i = 3
    while i * i <= x:
        if x % i == 0:
            divs[cnt] = i
            cnt += 1
            while x % i == 0:
                x //= i
        i += 2
    if x > 1:
        divs[cnt] = x
        cnt += 1
    g = 2
    while True:
        for i in range(cnt):
            if pow(g, (m - 1) // divs[i], m) == 1:
                break
        else:
            return g
        g += 1


def popcount(x):
    x = ((x >> 1) & 0x55555555) + (x & 0x55555555)
    x = ((x >> 2) & 0x33333333) + (x & 0x33333333)
    x = ((x >> 4) & 0x0F0F0F0F) + (x & 0x0F0F0F0F)
    x = ((x >> 8) & 0x00FF00FF) + (x & 0x00FF00FF)
    x = ((x >> 16) & 0x0000FFFF) + (x & 0x0000FFFF)
    return x


def tzcount(x):
    return popcount(~x & (x - 1))


def build_ntt():
    g = primitive_root(MOD)
    rank2 = tzcount(MOD - 1)
    root = [0] * (rank2 + 1)
    iroot = [0] * (rank2 + 1)
    rate2 = [0] * max(0, rank2 - 1)
    irate2 = [0] * max(0, rank2 - 1)
    rate3 = [0] * max(0, rank2 - 2)
    irate3 = [0] * max(0, rank2 - 2)
    root[rank2] = pow(g, (MOD - 1) >> rank2, MOD)
    iroot[rank2] = pow(root[rank2], MOD - 2, MOD)
    for i in range(rank2)[::-1]:
        root[i] = root[i + 1] * root[i + 1]
        root[i] %= MOD
        iroot[i] = iroot[i + 1] * iroot[i + 1]
        iroot[i] %= MOD
    prod = 1
    iprod = 1
    for i in range(rank2 - 1):
        rate2[i] = root[i + 2] * prod % MOD
        irate2[i] = iroot[i + 2] * iprod % MOD
        prod *= iroot[i + 2]
        prod %= MOD
        iprod *= root[i + 2]
        iprod %= MOD
    prod = 1
    iprod = 1
    for i in range(rank2 - 2):
        rate3[i] = root[i + 3] * prod % MOD
        irate3[i] = iroot[i + 3] * iprod % MOD
        prod *= iroot[i + 3]
        prod %= MOD
        iprod *= root[i + 3]
        iprod %= MOD
    return root, iroot, rate2, irate2, rate3, irate3


root, iroot, rate2, irate2, rate3, irate3 = build_ntt()
IMAG = root[2]
IIMAG = iroot[2]


def butterfly(a):
    n = len(a)
    h = (n - 1).bit_length()
    len_ = 0
    while len_ < h:
        if h - len_ == 1:
            p = 1 << (h - len_ - 1)
            rot = 1
            for s in range(1 << len_):
                offset = s << (h - len_)
                for i in range(p):
                    l = a[i + offset]
                    r = a[i + offset + p] * rot % MOD
                    a[i + offset] = (l + r) % MOD
                    a[i + offset + p] = (l - r) % MOD
                if s + 1 != 1 << len_:
                    rot *= rate2[(~s & -~s).bit_length() - 1]
                    rot %= MOD
            len_ += 1
        else:
            p = 1 << (h - len_ - 2)
            rot = 1
            for s in range(1 << len_):
                rot2 = rot * rot % MOD
                rot3 = rot2 * rot % MOD
                offset = s << (h - len_)
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
                if s + 1 != 1 << len_:
                    rot *= rate3[(~s & -~s).bit_length() - 1]
                    rot %= MOD
            len_ += 2


def butterfly_inv(a):
    n = len(a)
    h = (n - 1).bit_length()
    len_ = h
    while len_:
        if len_ == 1:
            p = 1 << (h - len_)
            irot = 1
            for s in range(1 << (len_ - 1)):
                offset = s << (h - len_ + 1)
                for i in range(p):
                    l = a[i + offset]
                    r = a[i + offset + p]
                    a[i + offset] = (l + r) % MOD
                    a[i + offset + p] = (l - r) * irot % MOD
                if s + 1 != (1 << (len_ - 1)):
                    irot *= irate2[(~s & -~s).bit_length() - 1]
                    irot %= MOD
            len_ -= 1
        else:
            p = 1 << (h - len_)
            irot = 1
            for s in range(1 << (len_ - 2)):
                irot2 = irot * irot % MOD
                irot3 = irot2 * irot % MOD
                offset = s << (h - len_ + 2)
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
                if s + 1 != (1 << (len_ - 2)):
                    irot *= irate3[(~s & -~s).bit_length() - 1]
                    irot %= MOD
            len_ -= 2


def convolution(a, b):
    n = len(a)
    m = len(b)
    if not n or not m:
        return []
    if min(n, m) <= 100:
        if n < m:
            n, m = m, n
            a, b = b, a
        res = [0] * (n + m - 1)
        for i in range(n):
            for j in range(m):
                res[i + j] += a[i] * b[j]
                res[i + j] %= MOD
        return res
    z = 1 << (n + m - 2).bit_length()
    a += [0] * (z - n)
    b += [0] * (z - m)
    butterfly(a)
    butterfly(b)
    for i in range(z):
        a[i] *= b[i]
        a[i] %= MOD
    butterfly_inv(a)
    a = a[: n + m - 1]
    iz = pow(z, MOD - 2, MOD)
    for i in range(n + m - 1):
        a[i] *= iz
        a[i] %= MOD
    return a


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    a = list(map(int, input().split()))
    b = list(map(int, input().split()))
    print(*convolution(a, b))
