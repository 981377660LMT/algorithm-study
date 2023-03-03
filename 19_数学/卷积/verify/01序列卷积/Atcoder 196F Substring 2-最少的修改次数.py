# 一维卷积应用- FFT求卷积(FFT在字符串匹配上的应用 https://blog.csdn.net/qq_31759205/article/details/80060918)
# 给定二进制字符串s和t
# !用最少的修改次数，使得S串中出现子串T(求两个字符串的最小汉明距离)
# 1<=len(t)<=len(s)<=1e6

# https://www.bilibili.com/read/cv10391549/

# !令f[i]表示S串中以i为起点的子串与T串相等字符的个数
# !则 f[i] = S0[i+k]*T0[k] + S1[i+k]*T1[k] (k=0,1,2,...,len(t)-1)
# !交叉项是多项式乘法,反转t串变为卷积的形式 fft求解
# !变为 f[i] = S0[p]*T0[q]+S1[p]*T1[q] (p+q=len(t)-1))
# 则答案为 len(t)-max(f[i]) (i=len(t)-1,len(t),...,len(s)-1)


from typing import Any
import numpy as np


def convolve(nums1: Any, nums2: Any) -> "np.ndarray":
    """fft求卷积"""
    n, m = len(nums1), len(nums2)
    ph = 1 << (n + m - 2).bit_length()
    T = np.fft.rfft(nums1, ph) * np.fft.rfft(nums2, ph)
    res = np.fft.irfft(T, ph)[: n + m - 1]
    return np.rint(res).astype(np.int64)


if __name__ == "__main__":

    s = input()
    t = input()
    t = t[::-1]
    n1, n2 = len(s), len(t)

    s0 = np.array([0] * n1, dtype=np.int64)
    s1 = np.array([0] * n1, dtype=np.int64)
    t0 = np.array([0] * n2, dtype=np.int64)
    t1 = np.array([0] * n2, dtype=np.int64)

    for i in range(n1):
        if s[i] == "0":
            s0[i] = 1
        else:
            s1[i] = 1

    for i in range(n2):
        if t[i] == "0":
            t0[i] = 1
        else:
            t1[i] = 1

    same0 = convolve(s0, t0)  # 相同字符0个数
    same1 = convolve(s1, t1)  # 相同字符1个数
    # print(cand1, cand2)
    res = 0
    for start in range(n2 - 1, n1):
        res = max(res, same0[start] + same1[start])
    print(n2 - res)
