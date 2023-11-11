# 小Q到商场购物，她有×元钱。已知商场有n个商品，每个商品有2种价格，
# 一种是原价商品，一种是打折价商品。打折价的商品虽然花钱较少，但由于是临期促销产品，
# 容易变质，因此小Q得到的幸福度不一定更高。
# 现在给出每个商品的原价×1、原价对应的幸福度w1、打折价x、打折价对应的幸福度w2。
# 小Q想知道，把正好×元钱花完，可以得到最大的幸福度是多少?
# 注︰每种商品最多只能买一件。如果买了打折的商品就不能再买原价。

# 如果小○能正好把x元钱花完，请输出她能获得的最大幸福度。否则直接输出。-1。
# n<=12

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    def dfs(index: int, maxHappy: int, curMoney: int) -> None:
        global res
        if curMoney == 0:
            res = max(res, maxHappy)
            return
        if index == n:
            return
        dfs(index + 1, maxHappy, curMoney)
        price1, happy1, price2, happy2 = goods[index]
        if price1 <= curMoney:
            dfs(index + 1, maxHappy + happy1, curMoney - price1)
        if price2 <= curMoney:
            dfs(index + 1, maxHappy + happy2, curMoney - price2)

    n, money = map(int, input().split())
    goods = [[] for _ in range(n)]
    for i in range(n):
        price1, happy1, price2, happy2 = map(int, input().split())
        goods[i] = [price1, happy1, price2, happy2]

    res = -1
    dfs(0, 0, money)
    print(res)

# 21
# 21 3 4 16 1 17 2 20 18 15 19 14 13 12 11 10 9 7 8 5 6
