# i*nums[i]

from itertools import accumulate
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n, m = map(int, input().split())
nums = list(map(int, input().split()))

preSum = [0] + list(accumulate(nums))

n = len(nums)
res, curSum = -INF, 0
for right in range(n):
    curSum += nums[right] * min(right + 1, m)
    if right >= m:
        curSum -= preSum[right] - preSum[right - m]
    if right >= m - 1:
        res = max(res, curSum)
print(res)
