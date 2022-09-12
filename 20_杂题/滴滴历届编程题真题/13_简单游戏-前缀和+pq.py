# 给出一个长度为 n 的数组 a，
# 你需要在这个数组中找到一个长度至少为 m 的区间，
# 使得这个区间内的数字的和尽可能小。

# !用pq维护最大和的子数组，当前的前缀和减左边最大前缀和即可
from heapq import heappush


n, m = [int(i) for i in input().split()]
nums = [int(i) for i in input().split()]

res = sum(nums[:m])
curSum = 0
leftSum = 0
pq = [0]

for i, cur in enumerate(nums):
    curSum += cur
    if i >= m:
        leftSum += nums[i - m]
        heappush(pq, -leftSum)
        leftMax = -pq[0]
        res = min(res, curSum - leftMax)

print(res)
