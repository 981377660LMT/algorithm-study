"""
题目大意:给出n个人的区间,表示这人在这个区间时间内会登录游戏,
现在请输出n个数k,表示有恰好k个人登录游戏的天数。
n<=2e5
ai,bi<=1e9

离散差分+扫描线
"""

from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    diff = defaultdict(int)
    for _ in range(n):
        left, length = map(int, input().split())
        right = left + length - 1
        diff[left] += 1
        diff[right + 1] -= 1

    res = [0] * (n + 1)
    curSum, pre = 0, -1
    for key in sorted(diff):
        delta = diff[key]
        if pre != -1:
            res[curSum] += key - pre
        curSum += delta
        pre = key
    print(*res[1:])
