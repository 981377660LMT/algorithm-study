# 所有(aj-ai)的和

from itertools import accumulate


n = int(input())
nums = list(map(int, input().split()))
preSum = [0] + list(accumulate(nums))

res = 0
for i in range(2, n + 1):
    res += (i - 1) * nums[i - 1] - preSum[i - 1]

print(res)
