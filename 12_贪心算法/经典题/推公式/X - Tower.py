# 给定 n 个方块，每个方块有三个属性：重量 w 、强度 s 与价值 v。
# 太郎想要选择价值总和尽可能大的一些方块摞在一起，
# !但是每个方块上方方块的总质量不能超过该方块的强度。
# 询问太郎能取得的最大价值。
# n<=1e3
# w s <=1e4
# v<=1e9

# 耍杂技的牛
# !Wi + Si
# !贪心排序+01背包(容量为s+w)
# 从上往下放
# 这个指标越小, 就应该放在上面, 越大, 就应该放在下面, 从而使最大风险减少.
# https://blog.51cto.com/u_15060513/3690362
# O(n*max(si))

# 为什么?:
# 偏序关系变为全序关系
# 假定现在已经堆叠了重量为 W 的方块，在下面要放下 i 与 j 两个方块，且最佳策略是 i 在上。
# 那么有 si<=W+wj sj>=W+wi 即si+wi<=sj+wj

import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
goods = []
for _ in range(n):
    weight, strong, value = map(int, input().split())
    goods.append((weight, strong, value))
goods.sort(key=lambda x: x[0] + x[1])  # !贪心排序

upper = int(2e4 + 10)
dp = [0] * upper
for w, s, v in goods:
    for cap in range(upper - 1, w - 1, -1):  # 当前所有物品的重量和
        if s >= cap - w:
            dp[cap] = max(dp[cap], dp[cap - w] + v)
print(max(dp))


# def dfs(index: int, preWeight: int) -> int:
#     if index == n:
#         return 0
#     hash_ = (index * int(1e4 + 1)) + preWeight
#     if memo[hash_] != -1:
#         return memo[hash_]
#     res = dfs(index + 1, preWeight)
#     w, s, v = blocks[index]
#     if s >= preWeight:
#         res = max(res, dfs(index + 1, preWeight + w) + v)
#     memo[hash_] = res
#     return res


# memo = [-1] * (int(1e4 + 1) * (n + 1))
# print(dfs(0, 0))
