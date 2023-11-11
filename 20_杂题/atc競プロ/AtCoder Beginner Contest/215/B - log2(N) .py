# 求满足2**k<=n 的最大的k
# n<=1e18

# 注意log2有精度问题
# !print(floor(log2(n)))  会有精度问题
# !用乘法而不是除法 用幂而不是对数

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    for i in range(63, -1, -1):
        if 2**i <= n:
            print(i)
            exit(0)
