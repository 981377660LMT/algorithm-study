# 比较pow(a,c)与pow(b,c)的大小
# -1e9<=a,b<=1e9
# 1<=c<=1e9
# pow(A,C)<pow(B,C) なら < を、pow(A,C)>pow(B,C) なら > を、pow(A,C)=pow(B,C) なら = を出力せよ。

# !技巧:看c模2等于0还是1
# !等于0则只需要比较2次幂 等于1则需要比较1次幂

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    a, b, c = map(int, input().split())

    mul = 1 if c & 1 else 2
    pow1, pow2 = pow(a, mul), pow(b, mul)
    if pow1 < pow2:
        print("<")
    elif pow1 > pow2:
        print(">")
    else:
        print("=")
