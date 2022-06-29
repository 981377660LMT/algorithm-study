# n个包裹 每个包裹里有两件物品
# 现在要从每个包裹里拿出一件物品 即n个物品凑齐s重量
# 输出具体方案 如果不能凑齐s重量 则输出Impossible
# n<=100
# s<=1e5

# 背包问题求具体方案的两类方法
# !1. bit压缩 这种 直接在dp里记录bit 比较慢
# !2. 倒序复原 同样适用于求最短路 根据dp数组倒退(dp数组本身包含很多信息)


# bit压缩
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n, s = map(int, input().split())  # n:物品数量 s:背包容量
goods = []
for _ in range(n):
    A, B = map(int, input().split())
    goods.append((A, B))
dp = [[False] * (s + 1) for _ in range(n + 1)]
dp[0][0] = True

for i, (a, b) in enumerate(goods):
    for cap in range(s + 1):
        if dp[i][cap]:
            if cap + a <= s:
                dp[i + 1][cap + a] = True
            if cap + b <= s:
                dp[i + 1][cap + b] = True

if not dp[-1][-1]:
    print("Impossible")
    exit(0)


# 倒退复原
res = []
cur = s
for i in range(n - 1, -1, -1):
    if cur - goods[i][0] >= 0 and dp[i][cur - goods[i][0]]:
        res.append("A")
        cur -= goods[i][0]
    else:
        res.append("B")
        cur -= goods[i][1]

print(''.join(res[::-1]))
