# 1≤i<j≤N を満たす整数の組 (i,j) のうち、
# 区間 i と区間 j が共通部分を持つようなものは幾つありますか？
# 1:[l,r]
# 2:[l,r)
# 3:(l,r]
# 4:(l,r)
# 开区间闭区间相交对数
# !技巧:开区间视为[left+EPS,right-EPS]
# !min(e1, e2) - max(s1, s2) >=0 就算有公共点

from itertools import combinations
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
EPS = 1e-6

if __name__ == "__main__":
    n = int(input())
    intervals = []
    for _ in range(n):
        kind, left, right = map(int, input().split())
        if kind == 1:
            intervals.append((left, right))
        elif kind == 2:
            intervals.append((left, right - EPS))
        elif kind == 3:
            intervals.append((left + EPS, right))
        else:
            intervals.append((left + EPS, right - EPS))

    res = 0
    for (s1, e1), (s2, e2) in combinations(intervals, 2):
        length = min(e1, e2) - max(s1, s2)
        res += length >= 0
    print(res)
