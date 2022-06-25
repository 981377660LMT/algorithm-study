# https://atcoder.jp/contests/practice2/submissions/31388245

MOD = 998244353
g = 3
ginv = 332748118
W = [pow(g, (MOD - 1) >> i, MOD) for i in range(24)]
Winv = [pow(ginv, (MOD - 1) >> i, MOD) for i in range(24)]


def fft(k, f):
    for l in range(k, 0, -1):
        d = 1 << l - 1
        U = [1]
        for i in range(d):
            U.append(U[-1] * W[l] % MOD)
        for i in range(1 << k - l):
            for j in range(d):
                s = i * 2 * d + j
                f[s], f[s + d] = (f[s] + f[s + d]) % MOD, U[j] * (f[s] - f[s + d]) % MOD


def fftinv(k, f):
    for l in range(1, k + 1):
        d = 1 << l - 1
        for i in range(1 << k - l):
            u = 1
            for j in range(i * 2 * d, (i * 2 + 1) * d):
                f[j + d] *= u
                f[j], f[j + d] = (f[j] + f[j + d]) % MOD, (f[j] - f[j + d]) % MOD
                u *= Winv[l]
                u %= MOD


def convolution(a, b):
    le = len(a) + len(b) - 1
    k = le.bit_length()
    n = 1 << k
    a = a + [0] * (n - len(a))
    b = b + [0] * (n - len(b))
    fft(k, a)
    fft(k, b)
    for i in range(n):
        a[i] *= b[i]
        a[i] %= MOD
    fftinv(k, a)

    ninv = pow(n, MOD - 2, MOD)
    for i in range(le):
        a[i] *= ninv
        a[i] %= MOD
    return a[:le]


n, m = map(int, input().split())
a = list(map(int, input().split()))
b = list(map(int, input().split()))

c = convolution(a, b)

print(*c)
