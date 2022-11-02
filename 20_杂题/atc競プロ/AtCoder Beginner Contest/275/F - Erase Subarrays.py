# 可以任意删除子数组
# !对x=1,2,...m 求出最少的删除次数 使得剩下的数组和为x
# !如果不能得到x,则输出-1
# n,m<=3000
# nums[i]<=3000

# !dp[index][curSum][preSelect] 表示前index个数,和为curSum,最后一个数是否被选中的最少删除次数

INF = int(1e18)

n, m = map(int, input().split())
nums = list(map(int, input().split()))
dp = [[INF] * 2 for _ in range(m + 1)]  # 0:不选 1:选
dp[0][1] = 0
for i in range(n):
    ndp = [[INF] * 2 for _ in range(m + 1)]
    for pre in range(m + 1):
        ndp[pre][0] = min(ndp[pre][0], dp[pre][0], dp[pre][1] + 1)  # 不选
        if pre + nums[i] <= m:
            ndp[pre + nums[i]][1] = min(ndp[pre + nums[i]][1], dp[pre][0], dp[pre][1])  # 选
    dp = ndp

for i in range(1, m + 1):
    res = min(*dp[i])
    print(res if res != INF else -1)
