# 题目大意:是你去买苹果，你有2种购买方式:
# (1)买一个苹果花费X元 (2)买3个苹果花费Y元。
# 然后问你最少用多少钱可以买到N个苹果

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


x, y, n = map(int, input().split())


# 那个更加划算
if 3 * x <= y:
    print(n * x)
    exit(0)

div, mod = divmod(n, 3)
print(div * y + mod * x)
