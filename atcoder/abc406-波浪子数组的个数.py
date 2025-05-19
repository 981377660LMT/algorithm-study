# C - ~ チルダ（tilde）
# abc406-波浪子数组的个数
# https://atcoder.jp/contests/abc406/tasks/abc406_c

from itertools import groupby
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    N = int(input())
    A = list(map(int, input().split()))
    diff = [b > a for a, b in zip(A, A[1:])]
    groups = [(v, len(list(group))) for v, group in groupby(diff)]
    res = 0
    for a, b, c in zip(groups, groups[1:], groups[2:]):
        if not b[0]:
            res += a[1] * c[1]
    print(res)
