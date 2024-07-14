import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# 1 以上
# N 以下の正整数
# x であって、ある正整数
# a と
# 2 以上の 正整数
# b を用いて
# x=a
# b
#   と表現できるものはいくつありますか？

# N<=1e18

# TODO:
# 1. 观察法
# 2. 因数容斥/因子容斥/因数莫比乌斯 ，加入容斥原理文件夹
if __name__ == "__main__":
    # 思路:
    # 1. 1以上N以下の数を全て素因数分解する
    # 2. 素因数分解の結果をもとに、a*bの形になる数を数える
    # 3. 1以上N以下の数から、a*bの形になる数を引く

    N = int(input())

    # 素因数分解
    # 1以上N以下の数を全て素因数分解する
