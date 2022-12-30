"""
Reference
https://github.com/atcoder/ac-library/blob/master/atcoder/convolution.hpp
https://github.com/atcoder/ac-library/blob/master/atcoder/internal_math.hpp
https://github.com/atcoder/ac-library/blob/master/document_en/convolution.md
https://github.com/atcoder/ac-library/blob/master/document_ja/convolution.md

注意此模板只能用于pypy3 
python3.8会超时,需要numpy的fft
"""
from typing import List


MOD = 998244353


def convolution(a: List[int], b: List[int]) -> List[int]:
    n = len(a)
    m = len(b)
    if not n or not m:
        return []
    if min(n, m) <= 60:
        return _convolution_naive(a, b)
    return _convolution_fft(a, b)


def _primitive_root(m):
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
        ok = True
        for i in range(cnt):
            if pow(g, (m - 1) // divs[i], m) == 1:
                ok = False
                break
        if ok:
            return g
        g += 1


class _FFTINFO:
    def __init__(self):
        self.g = _primitive_root(MOD)
        self.rank2 = ((MOD - 1) & (1 - MOD)).bit_length() - 1
        self.root = [0] * (self.rank2 + 1)
        self.root[self.rank2] = pow(self.g, (MOD - 1) >> self.rank2, MOD)
        self.iroot = [0] * (self.rank2 + 1)
        self.iroot[self.rank2] = pow(self.root[self.rank2], MOD - 2, MOD)
        for i in range(self.rank2 - 1, -1, -1):
            self.root[i] = self.root[i + 1] * self.root[i + 1] % MOD
            self.iroot[i] = self.iroot[i + 1] * self.iroot[i + 1] % MOD

        self.rate2 = [0] * max(0, self.rank2 - 1)
        self.irate2 = [0] * max(0, self.rank2 - 1)
        prod = 1
        iprod = 1
        for i in range(self.rank2 - 1):
            self.rate2[i] = self.root[i + 2] * prod % MOD
            self.irate2[i] = self.iroot[i + 2] * iprod % MOD
            prod *= self.iroot[i + 2]
            prod %= MOD
            iprod *= self.root[i + 2]
            iprod %= MOD

        self.rate3 = [0] * max(0, self.rank2 - 2)
        self.irate3 = [0] * max(0, self.rank2 - 2)
        prod = 1
        iprod = 1
        for i in range(self.rank2 - 2):
            self.rate3[i] = self.root[i + 3] * prod % MOD
            self.irate3[i] = self.iroot[i + 3] * iprod % MOD
            prod *= self.iroot[i + 3]
            prod %= MOD
            iprod *= self.root[i + 3]
            iprod %= MOD


info = _FFTINFO()


def _butterfly(a):
    n = len(a)
    h = (n - 1).bit_length()

    length = 0
    while length < h:
        if h - length == 1:
            p = 1 << (h - length - 1)
            rot = 1
            for s in range(1 << length):
                offset = s << (h - length)
                for i in range(p):
                    l = a[i + offset]
                    r = a[i + offset + p] * rot % MOD
                    a[i + offset] = (l + r) % MOD
                    a[i + offset + p] = (l - r) % MOD
                if s + 1 != (1 << length):
                    rot *= info.rate2[(~s & -~s).bit_length() - 1]
                    rot %= MOD
            length += 1
        else:
            # 4-base
            p = 1 << (h - length - 2)
            rot = 1
            imag = info.root[2]
            for s in range(1 << length):
                rot2 = rot * rot % MOD
                rot3 = rot2 * rot % MOD
                offset = s << (h - length)
                for i in range(p):
                    a0 = a[i + offset]
                    a1 = a[i + offset + p] * rot
                    a2 = a[i + offset + 2 * p] * rot2
                    a3 = a[i + offset + 3 * p] * rot3
                    a1na3imag = (a1 - a3) % MOD * imag
                    a[i + offset] = (a0 + a2 + a1 + a3) % MOD
                    a[i + offset + p] = (a0 + a2 - a1 - a3) % MOD
                    a[i + offset + 2 * p] = (a0 - a2 + a1na3imag) % MOD
                    a[i + offset + 3 * p] = (a0 - a2 - a1na3imag) % MOD
                if s + 1 != (1 << length):
                    rot *= info.rate3[(~s & -~s).bit_length() - 1]
                    rot %= MOD
            length += 2


def _butterfly_inv(a):
    n = len(a)
    h = (n - 1).bit_length()

    length = h  # a[i, i+(n<<length), i+2*(n>>length), ...] is transformed
    while length:
        if length == 1:
            p = 1 << (h - length)
            irot = 1
            for s in range(1 << (length - 1)):
                offset = s << (h - length + 1)
                for i in range(p):
                    l = a[i + offset]
                    r = a[i + offset + p]
                    a[i + offset] = (l + r) % MOD
                    a[i + offset + p] = (l - r) * irot % MOD
                if s + 1 != (1 << (length - 1)):
                    irot *= info.irate2[(~s & -~s).bit_length() - 1]
                    irot %= MOD
            length -= 1
        else:
            # 4-base
            p = 1 << (h - length)
            irot = 1
            iimag = info.iroot[2]
            for s in range(1 << (length - 2)):
                irot2 = irot * irot % MOD
                irot3 = irot2 * irot % MOD
                offset = s << (h - length + 2)
                for i in range(p):
                    a0 = a[i + offset]
                    a1 = a[i + offset + p]
                    a2 = a[i + offset + 2 * p]
                    a3 = a[i + offset + 3 * p]
                    a2na3iimag = (a2 - a3) * iimag % MOD
                    a[i + offset] = (a0 + a1 + a2 + a3) % MOD
                    a[i + offset + p] = (a0 - a1 + a2na3iimag) * irot % MOD
                    a[i + offset + 2 * p] = (a0 + a1 - a2 - a3) * irot2 % MOD
                    a[i + offset + 3 * p] = (a0 - a1 - a2na3iimag) * irot3 % MOD
                if s + 1 != (1 << (length - 2)):
                    irot *= info.irate3[(~s & -~s).bit_length() - 1]
                    irot %= MOD
            length -= 2


def _convolution_naive(a, b):
    n = len(a)
    m = len(b)
    ans = [0] * (n + m - 1)
    if n < m:
        for j in range(m):
            for i in range(n):
                ans[i + j] += a[i] * b[j]
                ans[i + j] %= MOD
    else:
        for i in range(n):
            for j in range(m):
                ans[i + j] += a[i] * b[j]
                ans[i + j] %= MOD
    return ans


def _convolution_fft(a, b):
    a = a.copy()
    b = b.copy()
    n = len(a)
    m = len(b)
    z = 1 << (n + m - 2).bit_length()
    a += [0] * (z - n)
    _butterfly(a)
    b += [0] * (z - m)
    _butterfly(b)
    for i in range(z):
        a[i] *= b[i]
        a[i] %= MOD
    _butterfly_inv(a)
    a = a[: n + m - 1]
    iz = pow(z, MOD - 2, MOD)
    for i in range(n + m - 1):
        a[i] *= iz
        a[i] %= MOD
    return a


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    a = list(map(int, input().split()))
    b = list(map(int, input().split()))
    print(*convolution(a, b))
