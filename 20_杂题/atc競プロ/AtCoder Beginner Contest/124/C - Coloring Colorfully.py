# 求最小翻转次数使得01串相邻两个元素不同
# !讨论第一个元素是否flip
# 101010...
# 010101...

import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    s = input()
    s1 = "01" * len(s)
    s2 = "10" * len(s)
    diff1 = sum(s[i] != s1[i] for i in range(len(s)))
    diff2 = sum(s[i] != s2[i] for i in range(len(s)))
    print(min(diff1, diff2))
