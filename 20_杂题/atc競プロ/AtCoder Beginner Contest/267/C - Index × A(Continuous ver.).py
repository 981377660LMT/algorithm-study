from itertools import accumulate, combinations
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n, m = map(int, input().split())
nums = list(map(int, input().split()))

# 定长滑窗 TODO
preSum = [0] + list(accumulate(nums))
res = -INF
curSum = 0
for i in range(n):
    curSum += nums[i] * min(i + 1, m)
    if i >= m:
        curSum -= preSum[i] - preSum[i - m]
    # print(curSum, nums[i])
    if i >= m - 1:
        res = max(res, curSum)
print(res)
