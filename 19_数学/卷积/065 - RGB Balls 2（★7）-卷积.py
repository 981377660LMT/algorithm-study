"""二维卷积"""

# 三种颜色的球 一共选K个
# R+G <= X
# G+B <= Y
# B+R <= Z
# 一共有多少种选法

# !条件转换,考虑反面:
# B>=K-X
# R>=K-Y
# G>=K-Z

# !卷积可以求解 畳み込み (画成二维矩阵 X+Y为定值K时 卷积值为Y=-X+K的对角线的和)
import sys
import numpy as np

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)
fac = [1]
ifac = [1]
for i in range(1, int(4e5) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


def _convolution(A, B):
    n, m = len(A), len(B)
    ph = 1 << (n + m - 2).bit_length()
    T = np.fft.rfft(A, ph) * np.fft.rfft(B, ph)
    res = np.fft.irfft(T, ph)[: n + m - 1]
    return np.rint(res).astype(np.int64)


def convolution(A: "np.ndarray", B: "np.ndarray", mod: int, s=10) -> "np.ndarray":
    s2 = s << 1
    mask = (1 << s) - 1

    m0, m1, m2 = A & mask, (A >> s) & mask, A >> s2
    n0, n1, n2 = B & mask, (B >> s) & mask, B >> s2

    p_0 = m0 + m2
    p0 = m0
    p1 = p_0 + m1
    pm1 = p_0 - m1
    pm2 = ((pm1 + m2) << 1) - m0
    pinf = m2

    q_0 = n0 + n2
    q0 = n0
    q1 = q_0 + n1
    qm1 = q_0 - n1
    qm2 = ((qm1 + n2) << 1) - n0
    qinf = n2

    r0 = _convolution(p0, q0)
    r1 = _convolution(p1, q1)
    rm1 = _convolution(pm1, qm1)
    rm2 = _convolution(pm2, qm2)
    rinf = _convolution(pinf, qinf)

    r_0 = r0
    r_4 = rinf
    r_3 = (rm2 - r1) // 3
    r_1 = (r1 - rm1) >> 1
    r_2 = rm1 - r0
    r_3 = ((r_2 - r_3) >> 1) + (rinf << 1)
    r_2 += r_1 - r_4
    r_1 -= r_3

    res = ((r_4 << s2) + (r_3 << s) + r_2) % mod
    return ((res << s2) + (r_1 << s) + r_0) % mod


R, G, B, K = map(int, input().split())
X, Y, Z = map(int, input().split())

nums1 = [0] * (B + 1)
for i in range(K - X, B + 1):
    nums1[i] = C(B, i)
nums2 = [0] * (R + 1)
for i in range(K - Y, R + 1):
    nums2[i] = C(R, i)
nums3 = [0] * (G + 1)
for i in range(K - Z, G + 1):
    nums3[i] = C(G, i)

nums1, nums2, nums3 = (
    np.array(nums1, dtype=np.int64),
    np.array(nums2, dtype=np.int64),
    np.array(nums3, dtype=np.int64),
)

print(convolution(nums1, convolution(nums2, nums3, MOD), MOD)[K])
