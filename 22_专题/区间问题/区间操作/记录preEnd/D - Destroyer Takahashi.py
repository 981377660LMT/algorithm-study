# 一次攻击可以击破长度为d的范围内的所有区间
# 击破所有区间至少需要多少次攻击
# n<=2e5 d<=1e9 lefti<=righti<=1e9

# !区间结束端点贪心排序
# 记录preEnd

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, d = map(int, input().split())
    intervals = []
    for _ in range(n):
        start, end = map(int, input().split())
        intervals.append((start, end))

    intervals.sort(key=lambda x: x[1])  # !按照区间结束端点排序
    res, preEnd = 0, -INF
    for curStart, curEnd in intervals:
        if curStart > preEnd:
            res += 1
            preEnd = curEnd + d - 1
    print(res)
