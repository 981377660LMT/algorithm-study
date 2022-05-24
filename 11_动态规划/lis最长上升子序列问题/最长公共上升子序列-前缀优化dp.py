# 最长公共上升子序列
# n<=3e3

# 定义dp数组：dp[i][j]表示第一个序列的前i个元素与第二个序列前j个字母以b[j]结尾的最长公共上升子序列长度
# 公共子序列既要最长、又要上升，上升需要额外的信息，所以规定dp[i,j]]以b[j]结尾，对]a[i]分析
# dp递推式：如果不选a[i]元素，最长公共上升子序列长度为dp[i-1][j];如果选a[i]元素,在a[i] == b[j]的前提，
# 内部按照以b[1]~b[j-1]结尾的上升子序列划分。
# 表达式为dp[i][j]=max(dp[i-1][1]+1,dp[i-1][2]+1,...dp[i-1][j-1]+1)
# O(n^3)
n = int(input())
nums1 = [0] + list(map(int, input().split()))
nums2 = [0] + list(map(int, input().split()))

dp = [[0 for _ in range(n + 1)] for _ in range(n + 1)]

for i in range(1, n + 1):
    for j in range(1, n + 1):
        dp[i][j] = dp[i - 1][j]  # 不选a[i]元素
        if nums1[i] == nums2[j]:  # 选择a[i]元素
            dp[i][j] = 1
            for k in range(1, j):
                if nums2[k] < nums2[j]:
                    dp[i][j] = max(dp[i][j], dp[i - 1][k] + 1)
res = 0
for i in range(1, n + 1):
    res = max(res, dp[n][i])
print(res)

# 因为每次转移都需要从前缀的状态转移过来，
# 所以可以用一个变量maxLen来储存当前a[i] > b[k]时dp[i-1][k]+1的最大值将复杂度降为二维。
# O(n^2)

n = int(input())
nums1 = [0] + list(map(int, input().split()))
nums2 = [0] + list(map(int, input().split()))

dp = [[0 for _ in range(n + 1)] for _ in range(n + 1)]

for i in range(1, n + 1):
    max_ = 1  # 记录最长的长度  表示 max{f[i-1][k] + 1, 0 <= k < j, b[k] < b[j]}
    for j in range(1, n + 1):
        dp[i][j] = dp[i - 1][j]  # 不选a[i]元素
        if nums1[i] == nums2[j]:  # 选择a[i]元素
            dp[i][j] = max(dp[i][j], max_)
        if nums1[i] > nums2[j]:
            max_ = max(max_, dp[i - 1][j] + 1)

res = 0
for i in range(1, n + 1):
    res = max(res, dp[n][i])
print(res)

