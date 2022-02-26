# 你打开了美了么外卖，选择了一家店，你手里有一张满 X 元减 10 元的券，
# 店里总共有 n 种菜，第 i 种菜一份需要 Ai 元，
# 因为你不想吃太多份同一种菜，
# 所以每种菜你最多只能点一份，
# 现在问你最少需要选择多少元的商品才能使用这张券。

# dp[I]是消费达到I 元所需的最低消费

n, X = map(int, input().split())
nums = list(map(int, input().split()))
dp = [int(1e20)] * (X + 1)
for i in range(n):
    for j in range(X, -1, -1):
        if j > nums[i]:
            dp[j] = min(dp[j], dp[j - nums[i]] + nums[i])
        else:
            # 注意这里
            dp[j] = min(dp[j], nums[i])

print(dp[X])

