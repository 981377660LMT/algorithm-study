# n张牌 所有牌正反面所有数的和为m
# 使得正面朝上的数和为k(k=0-m)时 至少需要翻动多少张牌 不能则返回-1
# (n,m<=2e5)

# 使用二进制优化的多重背包 两种极端情况
# 1. m个 0/1  此时复杂度为O(mlogn)
# 2. 0/1  0/2 0/3 0/4 .... sqrt(m) 个 此时复杂度为O(m*sqrt(m))

from collections import defaultdict, deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n, m = map(int, input().split())
    A, B = [], []
    for _ in range(n):
        a, b = map(int, input().split())
        A.append(a)
        B.append(b)

    # !1.普通的01背包dp O(nm) TLE
    # diff = [b - a for a, b in zip(A, B)]  # 每个背包的价值
    # dp = [INF] * (m + 1)
    # dp[sum(A)] = 0
    # for num in diff:
    #     ndp = dp[:]
    #     for pre in range(m + 1):
    #         if 0 <= pre + num <= m:
    #             ndp[pre + num] = min(ndp[pre + num], dp[pre] + 1)
    #     dp = ndp

    # for num in dp:
    #     print(num if num != INF else -1)

    # !2.二进制优化的多重背包 转换成0-1背包 AC
    diff = defaultdict(int)
    for a, b in zip(A, B):
        diff[b - a] += 1

    goods = []
    for cost, count in diff.items():
        size = 1
        while size <= count:
            goods.append((cost * size, size))
            count -= size
            size *= 2
        if count:
            goods.append((cost * count, count))

    dp = [INF] * (m + 1)
    dp[sum(A)] = 0
    for cost, count in goods:
        ndp = dp[:]
        for pre in range(m + 1):
            if 0 <= pre + cost <= m:
                ndp[pre + cost] = min(ndp[pre + cost], dp[pre] + count)
        dp = ndp

    for num in dp:
        print(num if num != INF else -1)

    # # !3.单调队列优化的多重背包 WA
    # diff = defaultdict(int)
    # for a, b in zip(A, B):
    #     diff[b - a] += 1

    # dp = [INF] * (m + 1)
    # dp[sum(A)] = 0
    # queue = deque()  # 单增的单调队列 队首最小值
    # for cost, count in diff.items():
    #     for j in range(cost):
    #         queue.clear()  # 注意每次新的循环都需要初始化队列
    #         remain = (m - j) // cost  # 剩余的空间最多还能放几个当前物品
    #         for k in range(remain + 1):
    #             val = dp[k * cost + j] - k * 1
    #             while queue and val <= queue[-1][1]:
    #                 queue.pop()
    #             queue.append((k, val))
    #             while queue and queue[0][0] < k - count:  # 存放的个数不能超出物品数量，否则弹出
    #                 queue.popleft()
    #             dp[k * cost + j] = queue[0][1] + k * 1

    # for num in dp:
    #     print(num if num != INF else -1)
