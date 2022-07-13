# 求1<=i,j<=n 的对数 (n<=1e5)
# 使得 ai < i < aj < j

# 前缀和

from typing import List


def solve(nums: List[int]) -> int:
    n = len(nums)
    preSum = [0]  # ai < i 的前缀和
    for i, num in enumerate(nums):
        preSum.append(preSum[-1] + int(num < i))
    res = 0
    for i in range(n):
        if nums[i] < i and nums[i] - 1 >= 0:
            res += preSum[nums[i] - 1]
    return res


print(solve([1, 1, 2, 3, 8, 2, 1, 4]))
