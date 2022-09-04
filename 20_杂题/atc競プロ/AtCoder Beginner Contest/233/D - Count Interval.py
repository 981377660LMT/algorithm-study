# 给一个数组nums和一个整数k,计算有多少个(L,R)满足数组A的[L,R]的和为K
# n<=2e5
# -1e9<=nums[i]<=1e9
# -1e15<=K<=1e15

# 哈希表+前缀和
from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, k = map(int, input().split())
nums = list(map(int, input().split()))
res = 0
preSum = defaultdict(int, {0: 1})
curSum = 0
for num in nums:
    curSum += num
    need = curSum - k
    res += preSum[need]
    preSum[curSum] += 1
print(res)
