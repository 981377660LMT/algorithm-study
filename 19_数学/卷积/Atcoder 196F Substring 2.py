# 一维卷积应用- FFT求卷积
# 给定二进制字符串s和t
# !用最少的修改次数，使得S串中出现子串T(求两个字符串的最小汉明距离)
# 1<=len(t)<=len(s)<=1e6

# https://www.bilibili.com/read/cv10391549/

# !令f[i]表示S串中以i为起点的子串与T串相等字符的个数
# !则 f[i] = S0[i+k]*T0[k] + S1[i+k]*T1[k] (k=0,1,2,...,len(t)-1)
# !交叉项是多项式乘法,反转t串变为卷积的形式 fft求解
# !变为 f[i] = S0[p]*T0[q]+S1[p]*T1[q] (p+q=len(t)-1))
# 则答案为 len(t)-max(f[i]) (i=len(t)-1,len(t),...,len(s)-1)

import sys
import numpy as np

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def convolve(a: np.ndarray, b: np.ndarray) -> np.ndarray:
    """fft求卷积"""
    fftLen = 1
    while 2 * fftLen < len(a) + len(b) - 1:
        fftLen *= 2
    fftLen *= 2
    Fa = np.fft.rfft(a, fftLen)
    Fb = np.fft.rfft(b, fftLen)
    Fc = Fa * Fb
    res = np.fft.irfft(Fc, fftLen)
    res = np.rint(res).astype(np.int64)
    return res[: len(a) + len(b) - 1]


if __name__ == "__main__":

    s = input()
    t = input()
    t = t[::-1]
    sLen, tLen = len(s), len(t)

    s0 = np.array([0] * sLen)
    s1 = np.array([0] * sLen)
    t0 = np.array([0] * tLen)
    t1 = np.array([0] * tLen)

    for i in range(sLen):
        if s[i] == "0":
            s0[i] = 1
        else:
            s1[i] = 1

    for i in range(tLen):
        if t[i] == "0":
            t0[i] = 1
        else:
            t1[i] = 1

    cand1 = convolve(s0, t0)  # 相同字符0个数
    cand2 = convolve(s1, t1)  # 相同字符1个数
    # print(cand1, cand2)
    res = 0
    for start in range(tLen - 1, sLen):
        res = max(res, cand1[start] + cand2[start])
    print(tLen - res)
