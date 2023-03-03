# https://judge.yosupo.jp/submission/45053

import numpy as np
import sys

# Karatsuba (爆速で、感動…)
# 行列積の精度がどの程度高いか不明なのでhackcaseが作れるんじゃないか説はある
# 快速矩阵乘法


def mat_prod(A, B, mod):
    A0 = (A & (1 << 15) - 1).astype(np.float64)
    A1 = (A >> 15).astype(np.float64)
    B0 = (B & (1 << 15) - 1).astype(np.float64)
    B1 = (B >> 15).astype(np.float64)
    C0 = A0 @ B0
    C2 = A1 @ B1
    C1 = C2 + C0 - (A1 - A0) @ (B1 - B0)
    C = C0.astype(np.int64)
    C += C1.astype(np.int64) % mod * (1 << 15)
    C += C2.astype(np.int64) % mod * (1 << 30)
    return C % mod


N, M, K = map(int, input().split())
a = list(map(int, sys.stdin.read().split()))
A = np.array(a[: N * M], dtype=np.int32).reshape(N, M)
B = np.array(a[N * M :], dtype=np.int32).reshape(M, K)
C = mat_prod(A, B, 998244353).tolist()
for row in C:
    print(" ".join(map(str, row)))
