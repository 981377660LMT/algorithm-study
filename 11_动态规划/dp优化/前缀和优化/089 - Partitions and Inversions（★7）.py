"""
1. 滑窗预处理边界
2. 前缀和优化dp
"""
# !求数组的分割方案数
# !使得每一段的逆序对(冒泡排序交换次数) <= k
# n<=2e5
# 1<=nums[i]<=1e9

# !怎么dp 普通的dp求和是O(n^2)的
# !前缀和优化dp O(n)

from collections import defaultdict
import sys


sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


class BIT1:
    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] += delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


n, k = map(int, input().split())
nums = list(map(int, input().split()))
bit = BIT1(max(nums) + 10)

# !1.对于每一个right 双指针求允许的left边界 (指这一段内逆序对<=k) 尺取法
leftMosts = [0] * n
left = 0
invCount = 0
for right in range(n):
    bit.add(nums[right], 1)
    invCount += bit.queryRange(nums[right] + 1, bit.size)
    while left < right and invCount > k:
        bit.add(nums[left], -1)
        invCount -= bit.query(nums[left] - 1)
        left += 1
    leftMosts[right] = left


# !2.前缀和优化dp dp[i]=dp[leftMost[i]]...+dp[i-1]
# !求出dp[i]后再更新dpSum[i]
dp = defaultdict(int, {0: 1})
dpSum = defaultdict(int, {0: 1})
for i in range(1, n + 1):
    dp[i] = (dpSum[i - 1] - dpSum[leftMosts[i - 1] - 1]) % MOD
    dpSum[i] = (dpSum[i - 1] + dp[i]) % MOD
print(dp[n])
