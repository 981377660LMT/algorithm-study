# 1≤N≤3000
# dp[i][j]代表所有a[1 ~ i]和b[1 ~ j]中以nums2[j]结尾的公共上升子序列的集合；
# dp[i][j]的值等于该集合的子序列中长度的最大值；

n = int(input())
nums1 = [0] + list(map(int, input().split()))
nums2 = [0] + list(map(int, input().split()))

dp = [[0] * (n + 1) for _ in range(n + 1)]

for i in range(1, n + 1):
    maxv = 1  # 用maxv记录最长的长度
    for j in range(1, n + 1):
        dp[i][j] = dp[i - 1][j]  # 不选a[i]元素
        # 公共子序列
        if nums1[i] == nums2[j]:  # 选择a[i]元素
            dp[i][j] = max(dp[i][j], maxv)
        if nums1[i] > nums2[j]:
            maxv = max(maxv, dp[i - 1][j] + 1)
res = 0
for i in range(1, n + 1):
    res = max(res, dp[n][i])
print(res)

