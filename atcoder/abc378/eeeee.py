import sys

sys.setrecursionlimit(1 << 25)
input = sys.stdin.readline

MOD = 998244353  # 适用于 NTT 的质数
ROOT = 3


def bit_reverse(a):
    n = len(a)
    j = 0
    for i in range(1, n):
        bit = n >> 1
        while j & bit:
            j ^= bit
            bit >>= 1
        j ^= bit

        if i < j:
            a[i], a[j] = a[j], a[i]


def ntt(a, invert):
    n = len(a)
    bit_reverse(a)

    length = 2
    while length <= n:
        wlen = pow(ROOT, (MOD - 1) // length, MOD)
        if invert:
            wlen = pow(wlen, MOD - 2, MOD)
        for i in range(0, n, length):
            w = 1
            for j in range(i, i + length // 2):
                u = a[j]
                v = a[j + length // 2] * w % MOD
                a[j] = (u + v) % MOD
                a[j + length // 2] = (u - v + MOD) % MOD
                w = w * wlen % MOD
        length <<= 1

    if invert:
        inv_n = pow(n, MOD - 2, MOD)
        for i in range(n):
            a[i] = a[i] * inv_n % MOD


def multiply(a, b):
    n = 1
    while n < len(a) + len(b):
        n <<= 1
    a.extend([0] * (n - len(a)))
    b.extend([0] * (n - len(b)))
    ntt(a, False)
    ntt(b, False)
    for i in range(n):
        a[i] = a[i] * b[i] % MOD
    ntt(a, True)
    return a


if __name__ == "__main__":
    N, M = map(int, input().split())
    A = list(map(int, input().split()))

    prefix_mod = [0] * (N + 1)
    for i in range(N):
        prefix_mod[i + 1] = (prefix_mod[i] + A[i]) % M

    freq = [0] * M
    for val in prefix_mod:
        freq[val] += 1

    # 扩展频率数组用于卷积
    size = 1
    while size < 2 * M:
        size <<= 1

    freq_a = freq + [0] * (size - M)
    freq_b = freq[::-1] + [0] * (size - M)

    # 使用 NTT 进行卷积
    conv = multiply(freq_a, freq_b)

    # 计算总和 S
    S = 0
    for d in range(M):
        count = conv[M - 1 + d]
        S = (S + count * d) % MOD

    print(S)
