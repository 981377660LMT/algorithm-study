# 累乘和 (幂级数求和/等比数列求和)
# 1<=q,n,MOD<=1e9
# https://atcoder.jp/contests/abc293/editorial/5955

# 令 a[n]=∑q^i(i=0,1,2,...,n-1)
# 则 a[0]=0,a[n+1]=q*a[n] + 1
# !可以写成矩阵快速幂的形式
# [an]  =  [q 1] ^ n * [0]
# [1 ]     [0 1]       [1]
# !等差乘等比数列同理

import numpy as np


def matqpow2(base: "np.ndarray", exp: int, mod: int) -> "np.ndarray":
    """np矩阵快速幂"""
    res = np.eye(*base.shape, dtype=np.uint64)
    while exp:
        if exp & 1:
            res = (res @ base) % mod
        base = (base @ base) % mod
        exp >>= 1
    return res


def powerSum(q: int, n: int, mod: int) -> int:
    """q^0 + q^1 + ... + q^(n-1) mod mod"""
    base = np.array([[q, 1], [0, 1]], dtype=np.uint64)
    trans = matqpow2(base, n, mod)
    res = trans @ np.array([[0], [1]], dtype=np.uint64)
    return res[0][0]


if __name__ == "__main__":
    q, n, MOD = map(int, input().split())
    print(powerSum(q, n, MOD))
