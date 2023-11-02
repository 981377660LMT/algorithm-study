# https://atcoder.jp/contests/dp/tasks/dp_e
# 超大容量01背包 -> 维度转换


from typing import List

INF = int(1e18)


# n<=100
# values[i]<=1e3
# weights[i]<=1e9
# maxCapacity <=1e9
# dp[i][j] 表示前i个物品，价值为j时最小的重量.
def knapsack2(values: List[int], weights: List[int], maxCapacity: int) -> int:
    def min2(a: int, b: int) -> int:
        return a if a < b else b

    valueSum = sum(values)
    dp = [INF] * (valueSum + 1)
    dp[0] = 0
    for v, w in zip(values, weights):
        for j in range(valueSum, v - 1, -1):
            dp[j] = min2(dp[j], dp[j - v] + w)

    res = 0
    for i in range(valueSum, -1, -1):
        if dp[i] <= maxCapacity:
            res = i
            break
    return res


if __name__ == "__main__":
    n, maxCapacity = map(int, input().split())
    values = []
    weights = []
    for _ in range(n):
        w, v = map(int, input().split())
        values.append(v)
        weights.append(w)
    print(knapsack2(values, weights, maxCapacity))
