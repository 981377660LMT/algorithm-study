# https://leetcode.cn/problems/VdG6tT/solution/si-bian-xing-bu-deng-shi-you-hua-dp-by-k-bp1t/

# N个工程师分成 K 个组(N<=1e4,K<=1e2) 组内噪声
# noise(l, r) = sum(A[l], A[l + 1], ..., A[r]) * (r - l + 1)
# 最小化总噪声

# 记dp[i][k]为前i个人分k组的最小噪声值，
# 则dp[i][k] = min(dp[j][k-1]+(preSum[i]-preSum[j])*(i-j+1))
# 直接猜测可以用四边形不等式


# !暴力dp O(n^2*k)
from itertools import accumulate

n, k = map(int, input().split())
nums = list(map(int, input().split()))
preSum = list(accumulate(nums, initial=0))
dp = [[int(1e20)] * (k + 5) for _ in range(n + 5)]

for i in range(1, n + 1):
    for g in range(1, min(k, i) + 1):
        if g == 1:
            dp[i][g] = preSum[i] * i
        else:
            # !真的需要遍历这么多吗 其实不需要
            # !因为后面pi里i更近(这一段人更少) 但是费用却比之前少 那么之前的位置就可以退役了
            for pi in range(g - 1, i):
                cand = dp[pi][g - 1] + (preSum[i] - preSum[pi]) * (i - pi)
                if cand < dp[i][g]:
                    dp[i][g] = cand

prePos = [g - 1 for g in range(1, k + 5)]
for i in range(1, n + 1):
    for g in range(1, min(k, i) + 1):
        if g == 1:
            dp[i][g] = preSum[i] * i
        else:
            # !对每一个固定的分组g 真的需要遍历这么多pi吗 其实不需要
            # !因为后面pi里i更近(这一段人更少) 但是费用却比之前少 那么之前的位置就可以退役了
            for pi in range(prePos[g], i):
                cand = dp[pi][g - 1] + (preSum[i] - preSum[pi]) * (i - pi)
                if cand < dp[i][g]:
                    dp[i][g] = cand
                    prePos[g] = pi

print(dp[n][k])

