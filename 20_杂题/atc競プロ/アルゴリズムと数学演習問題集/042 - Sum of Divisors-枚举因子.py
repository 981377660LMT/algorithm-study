# 将正整数X的约数个数记作 f(X)
# 求 ∑i*f(i) 的值  i=1,2,...,N
# N<=1e7
# !すべての整数で約数を見てるとO(n**2)かかる
# !n_iを約数に持つものを軸に考えると
# !初項n_i、項数n//n_i、公差n_iの等差数列になる
# !これはO(1)で求まるので
# !最終的にO(n)で答えが出る

import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n = int(input())

res = 0
for f in range(1, n + 1):
    # 每个因子的贡献 在哪里出现
    first, last, count = f, (n // f) * f, n // f
    res += (first + last) * count // 2
print(res)
