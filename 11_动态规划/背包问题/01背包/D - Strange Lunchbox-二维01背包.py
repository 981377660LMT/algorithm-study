# 奇怪的便当(lunchbox)
# 二维01背包
# 有n个物品 价值1为a 价值2为b
# !选出一些物品使得价值1的和>=x 价值2的和>=y
# !问最少需要选出多少个物品

# n,x,y,a,b<=300
# !注意为了减少状态数 超出x的数一律表示为x

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    x, y = map(int, input().split())
    goods = [tuple(map(int, input().split())) for _ in range(n)]

    # @lru_cache(None)
    # def dfs(index: int, countA: int, countB: int) -> int:
    #     if countA >= x and countB >= y:
    #         return 0
    #     if index == n:
    #         return INF

    #     res = dfs(index + 1, countA, countB)  # 不选
    #     a, b = goods[index]
    #     res = min(res, dfs(index + 1, min(x, countA + a), min(y, countB + b)) + 1)  # 选
    #     return res

    # res = dfs(0, 0, 0)
    # dfs.cache_clear()
    # print(res if res != INF else -1)

    dp = [[INF] * (y + 1) for _ in range(x + 1)]
    dp[0][0] = 0
    for a, b in goods:
        ndp = [[INF] * (y + 1) for _ in range(x + 1)]
        for pre1 in range(x + 1):
            for pre2 in range(y + 1):
                # 不选
                ndp[pre1][pre2] = min(ndp[pre1][pre2], dp[pre1][pre2])
                # 选
                cur1, cur2 = min(x, pre1 + a), min(y, pre2 + b)
                ndp[cur1][cur2] = min(ndp[cur1][cur2], dp[pre1][pre2] + 1)
        dp = ndp

    res = dp[x][y]
    print(res if res != INF else -1)
