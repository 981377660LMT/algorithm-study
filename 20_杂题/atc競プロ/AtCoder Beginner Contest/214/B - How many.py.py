# a+b+c<=S 且 a*b*c<=T
# !求非负整数三元组(a,b,c)的个数
# S,T<=1e18

from itertools import product
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
