from typing import List

from typing import List


MOD = int(1e9) + 7
INV = pow(2, MOD - 2, MOD)


def xorConvolution(a: List[int], b: List[int]) -> List[int]:
    a, b = a[:], b[:]
    a = _walshHadamardTransform(a, 1)
    b = _walshHadamardTransform(b, 1)
    for i in range(len(a)):
        a[i] = a[i] * b[i] % MOD
    res = _walshHadamardTransform(a, INV)
    return res


def _walshHadamardTransform(f, op):
    n = len(f)
    l_, k = 2, 1
    while l_ <= n:
        for i in range(0, n, l_):
            for j in range(k):
                f[i + j], f[i + j + k] = (f[i + j] + f[i + j + k]) * op % MOD, (
                    f[i + j] + MOD - f[i + j + k]
                ) * op % MOD
        l_, k = l_ << 1, k << 1
    return f


if __name__ == "__main__":
    # https://judge.yosupo.jp/submission/129301
    MOD = 998244353
    INV = (MOD + 1) // 2
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))
    conv = xorConvolution(nums1, nums2)
    print(*conv)
