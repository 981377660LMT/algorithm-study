# 区间内最大公约数不等于1的对数
# 给你一个区间[L,R]，
# 求有多少个pair(x,y)满足条件
# !gcd(x,y)!=1 and x//gcd(x,y) !=1 and y//gcd(x,y) !=1
# 1<=left<=right<=1e6

# !枚举因子+容斥原理(减去x或y正好等于gcd的对数)
# https://atcoder.jp/contests/abc206/tasks/abc206_e

from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    MAX = int(1e6) + 10
    left, right = map(int, input().split())
    ...


# TODO
