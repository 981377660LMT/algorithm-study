from typing import List

MOD = 998244353


def fwt_xor(arr: List[int], inv=False) -> None:
    n = len(arr)
    assert (n & (n - 1)) == 0, "Length of arr must be a power of 2"
    for i in range(n.bit_length()):
        m = b = 1 << i
        while m < n:
            a0, a1 = arr[m ^ b], arr[m]
            arr[m ^ b], arr[m] = a0 + a1, a0 - a1
            m = m + 1 | b
    if inv:
        inv_ = pow(n, -1, MOD)
        for i in range(n):
            arr[i] = (arr[i] * inv_) % MOD


def fwt_or(arr: List[int], inv=False) -> None:
    if not inv:
        _subset_zeta_transform(arr)
    else:
        _subset_moebius_transform(arr)


def fwt_and(arr: List[int], inv=False) -> None:
    if not inv:
        _superset_zeta_transform(arr)
    else:
        _superset_moebius_transform(arr)


def _subset_zeta_transform(arr: List[int]) -> None:
    n = len(arr)
    assert (n & (n - 1)) == 0, "Length of arr must be a power of 2"
    for i in range(n.bit_length()):
        m = b = 1 << i
        while m < n:
            arr[m] += arr[m ^ b]
            m = m + 1 | b
    for i in range(n):
        arr[i] %= MOD


def _subset_moebius_transform(arr: List[int]) -> None:
    n = len(arr)
    assert (n & (n - 1)) == 0, "Length of arr must be a power of 2"
    for i in range(n.bit_length()):
        m = b = 1 << i
        while m < n:
            arr[m] -= arr[m ^ b]
            m = m + 1 | b
    for i in range(n):
        arr[i] %= MOD


def _superset_zeta_transform(arr: List[int]) -> None:
    n = len(arr)
    assert (n & (n - 1)) == 0, "Length of arr must be a power of 2"
    for i in range(n.bit_length()):
        m = b = 1 << i
        while m < n:
            arr[m ^ b] += arr[m]
            m = m + 1 | b
    for i in range(n):
        arr[i] %= MOD


def _superset_moebius_transform(arr: List[int]) -> None:
    n = len(arr)
    assert (n & (n - 1)) == 0, "Length of arr must be a power of 2"
    for i in range(n.bit_length()):
        m = b = 1 << i
        while m < n:
            arr[m ^ b] -= arr[m]
            m = m + 1 | b
    for i in range(n):
        arr[i] %= MOD


if __name__ == "__main__":
    # https://www.luogu.com.cn/problem/P4717

    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))

    def solve(fwtFunc) -> None:
        a = nums1[:]
        b = nums2[:]
        fwtFunc(a, False)
        fwtFunc(b, False)
        for i in range(len(b)):
            a[i] = (a[i] * b[i]) % MOD
        fwtFunc(a, True)
        print(*a)

    solve(fwt_or)
    solve(fwt_and)
    solve(fwt_xor)
