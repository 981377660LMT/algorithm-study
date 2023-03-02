"""numpy fft 任意模数"""


import numpy as np


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


if __name__ == "__main__":
    MOD = int(1e9 + 7)
    n, m = map(int, input().split())
    A = np.array(list(map(int, input().split())), dtype=np.int64)
    B = np.array(list(map(int, input().split())), dtype=np.int64)
    C = convolution(A, B, MOD)
    print(*C)
