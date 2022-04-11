from typing import List
from heapq import heapify, heapreplace

MOD = int(1e9 + 7)
INF = int(1e20)

# 6039. K 次增加后的最大乘积

# x>1时 给x加1后 `总体增大了 (x+1)/x 倍` 所以x越小越好
# x=0时 也是一样

# 1 <= nums.length, k <= 1e5

# 如果k很大 需要考虑二分
class Solution:
    def maximumProduct(self, nums: List[int], k: int) -> int:
        heapify(nums)
        for _ in range(k):
            heapreplace(nums, nums[0] + 1)
        res = 1
        for num in nums:
            res *= num
            res %= MOD
        return res % MOD


print(Solution().maximumProduct(nums=[0, 4], k=5))
