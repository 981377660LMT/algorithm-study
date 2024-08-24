# abc368-e-trainDelay
# 电车延误(高铁晚点)
# https://atcoder.jp/contests/abc368/tasks/abc368_e
# 有n辆电车，每辆从A[i]开往B[i]，发车时间为S[i]，到达时间为T[i].
# 现在电车0需要延误X0分钟，对电车1-n，求最小延误时间,使得X1+X2...+Xn-1最小.
# 延误时间需要满足：
# 对所有的二元组(i,j)，如果B[i]==A[j]且T[i]<=S[j]，则需要有 T[i]+X[i]<=S[j]+X[j]
# 也即，原来可以换乘的电车，延误之后也可以换乘

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    N, M, X0 = map(int, input().split())
    A, B, S, T = [0] * M, [0] * M, [0] * M, [0] * M
    for i in range(M):
        a, b, s, t = map(int, input().split())
        A[i], B[i], S[i], T[i] = a - 1, b - 1, s, t
