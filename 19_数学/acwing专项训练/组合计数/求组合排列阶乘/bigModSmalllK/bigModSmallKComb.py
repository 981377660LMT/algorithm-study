def bigModSmallKComb(n: int, k: int, mod: int) -> int:
    """求组合数.适用于 n 巨大但 k 或 n-k 较小的情况."""
    if k > n - k:
        k = n - k
    a, b = 1, 1
    for i in range(1, k + 1):
        a = a * n % mod
        n -= 1
        b = b * i % mod
    return a * pow(b, mod - 2, mod) % mod


if __name__ == "__main__":
    print(bigModSmallKComb(100, 20, int(1e9 + 7)))
