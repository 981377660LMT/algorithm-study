# n = int(input())
# nums = list(map(int, input().split()))
# 数组划分成三个子集最大值最小-三组石头

n = 2
nums = [2, 6]
if n <= 3:
    print(max(nums))
else:
    # 普通的01背包 前i个物品 容量为j时 最多可以放多少物品
    total = sum(nums)
    dp = [[0] * (total + 1) for _ in range(n + 1)]
    for i in range(1, n + 1):
        for j in range(1, total + 1):
            if j >= nums[i - 1]:
                dp[i][j] = max(dp[i - 1][j], dp[i - 1][j - nums[i - 1]] + nums[i - 1])
            else:
                dp[i][j] = dp[i - 1][j]

    res = int(1e20)
    for i in range(1, total + 1):
        a = dp[-1][i]
        b = dp[-1][(total - a) // 2]
        c = total - a - b
        res = min(res, max(a, b, c))
    print(res)
