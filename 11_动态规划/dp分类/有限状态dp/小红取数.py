# 小红拿到了一个数组，她想取一些数使得取的数之和尽可能大，但要求这个和必须是 k 的倍数。
# dp[index][mod] 采用滚动数组更新


n, k = list(map(int, input().split()))
nums = list(map(int, input().split()))
# n, k = 5, 5
# nums = [4, 8, 2, 9, 1]

INF = int(1e100)  # 要取得足够大
dp = [-INF] * k
dp[0] = 0


for num in nums:
    mod = num % k
    copy = dp[:]
    for i in range(k):
        dp[i] = max(dp[i], num + copy[(i - mod)])

print(dp[0] if dp[0] != 0 else -1)
