# F - Knapsack for All Segments
# 区间的01背包
# https://atcoder.jp/contests/abc159/tasks/abc159_f

# 输入 n(1≤n≤3000) s(1≤s≤3000) 和长为 n 的数组 a(1≤a[i]≤3000)。
# 定义 f(L,R) 等于：在子数组 a[L],a[L+1],...,a[R] 中，元素和恰好等于 s 的子序列的个数。
# 输出所有 f(L,R) 的和，其中 0≤L≤R<n。
# 模 998244353。

# !dp[i][sum] 表示子数组右端点为 i 时，子序列和为 sum 的方案数。
# !要求的就是dp[0][target]+...+dp[n-1][target]。
# jump or not jump =>
# dp[i][sum] = dp[i-1][sum-a[i]] + dp[i-1][sum]


from typing import List

MOD = 998244353


def knapsackForAllSegments(nums: List[int], target: int) -> int:
    dp = [0] * (target + 1)
    dp[0] = 1
    res = 0
    for i, num in enumerate(nums):
        ndp = dp[:]  # 不选
        ndp[0] = i + 2
        for s in range(target, num - 1, -1):
            ndp[s] += dp[s - num]
            ndp[s] %= MOD
        dp = ndp
        res += dp[target]
        res %= MOD
    return res


if __name__ == "__main__":
    n, target = map(int, input().split())
    nums = list(map(int, input().split()))
    print(knapsackForAllSegments(nums, target))
