# E - Square Price (价值为平方的分组背包/完全背包)
# https://www.cnblogs.com/zhr0102/p/18679014
# https://atcoder.jp/contests/abc389/tasks/abc389_e
# 题意：有n个商品，价格为pi，如果你买k件要花费k^2*pi，你有m块，问最多买多少件物品。
# n<=2e5, m<=1e18, pi<=2e9
#
# 题目等价于：
# 对于一个物品，在我们买了k−1件的前提下，购买第k件需要(k2−(k−1)2))pi=(2×k−1)pi元。
# 我们可以按照从便宜到贵的顺序尽可能多地购买物品，直到无法购买为止。
# !那么我们可以二分一个价格，把低于这个价格的都买上，剩下的钱尽可能买x+1元的物品。因为我们二分的是最大的价格，所以价值为x+1一定存在，而且我买不完。


from typing import List, Tuple


INF = int(1e20)


def solve(n: int, m: int, prices: list[int]) -> int:
    def check(mid: int) -> Tuple[bool, List[int]]:
        """
        购买所有价格<=mid的物品, 是否花费不超过m.
        物品的价格为 pi, 3*pi, 5*pi, ...
        """
        cost = 0
        curBuy = [0] * n
        for i, p in enumerate(prices):
            # (2 * k - 1) * p <= mid
            k = mid // p
            k += 1
            k //= 2
            curBuy[i] = k
            cost += k * k * p
            if cost > m:
                return False, curBuy
        return cost <= m, curBuy

    left, right = 1, INF
    buy = [0] * n
    while left <= right:
        mid = (left + right) // 2
        ok, curBuy = check(mid)
        if ok:
            buy = curBuy
            left = mid + 1
        else:
            right = mid - 1

    for i in range(n):
        m -= buy[i] * buy[i] * prices[i]
    for i in range(n):
        curPrice = (2 * buy[i] + 1) * prices[i]
        if m >= curPrice:
            buy[i] += 1
            m -= curPrice
    return sum(buy)


if __name__ == "__main__":
    n, m = map(int, input().split())
    prices = list(map(int, input().split()))
    print(solve(n, m, prices))
