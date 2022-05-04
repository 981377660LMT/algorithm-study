# 给定长度为n的数组，以及数字k，把数组分成连续的k部分，
# 每一份的分数等于子数组的和乘上长度，求数组总分最小值

# 样本输入
# 4 2 1 3 2 4 样本输出 20

# n <= 1e4  k<=min(n,100)
from itertools import accumulate


n, k = map(int, input().split())
nums = list(map(int, input().split()))
preSum = [0] + list(accumulate(nums))
dp = [[int(1e20)] * (k + 1) for _ in range(n + 1)]
dp[0][0] = 0
bestJ = list(range(k + 1))

# dp[i] = min(dp[j] + (i - j) * (p[i] - p[j])) # j < i 分成任意组

for i in range(1, n + 1):
    for group in range(1, min(k, i) + 1):
        if group == 1:
            dp[i][1] = preSum[i] * i
        else:
            # 可不可以不从这么多j转移过来呢:剪枝，记录之前最好的位置
            # (preSum[i] - preSum[j]) * (i - j) 是随长度递增的
            for j in range(bestJ[group], i):
                cand = dp[j][group - 1] + (preSum[i] - preSum[j]) * (i - j)
                if dp[i][group] > cand:
                    dp[i][group] = cand
                    bestJ[group] = j

print(dp[-1][-1])
# wqs 二分适用的题目类型：
