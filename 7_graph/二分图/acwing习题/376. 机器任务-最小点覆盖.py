from hungarian import Hungarian

import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# https://blog.csdn.net/baidu_23955875/article/details/46573537
# todo
# 简化版的匈牙利学习 （不判断二分图）
# 有两台机器 A，B 以及 K 个任务。
# 机器 A 有 N 种不同的模式（模式 0∼N−1），机器 B 有 M 种不同的模式（模式 0∼M−1）。
# 两台机器最开始都处于模式 0。
# 每个任务既可以在 A 上执行，也可以在 B 上执行。
# 对于每个任务 i，给定两个整数 a[i] 和 b[i]，表示如果该任务在 A 上执行，需要设置模式为 a[i]，如果在 B 上执行，需要模式为 b[i]。
# 任务可以以任意顺序被执行，但每台机器转换一次模式就要重启一次。
# 求怎样分配任务并合理安排顺序，能使机器重启次数最少。
# N,M<=100 k<1000

# !(最小点覆盖 = 最大匹配)


def solve() -> None:
    H = Hungarian(n * m, n * m)

    for _ in range(k):
        _, a, b = map(int, input().split())
        if a == 0 or b == 0:  # !初始的模式 不用切换
            continue
        H.addEdge(a, b)
    print(H.work())


while True:
    n, *rest = map(int, input().split())
    if len(rest) == 0:
        exit(0)
    m, k = rest[0], rest[1]
    solve()
