from typing import List

# 计算F(0), F(1), ..., F(n-1)中的最大值。
#  f(i)          = 0 * A[0] + 1 * A[1] + 2 * A[2] + .... +  (k-1) * A[k-1] + k * A[k]
#  f(i+1)        = 1 * A[0] + 2 * A[1] + 3 * A[2] + ...  +     k  * A[k-1] + 0 * A[k]

# => f(i+1) = f(i) + sum(Array) -  last element * array.length
class Solution:
    def maxRotateFunction(self, nums: List[int]) -> int:
        n = len(nums)
        total = sum(nums)
        res = sum(i * j for i, j in enumerate(nums))
        curSum = res

        for i in range(1, n):
            curSum += total - n * nums[-i]
            res = max(res, curSum)
        return res


print(Solution().maxRotateFunction([4, 3, 2, 6]))
# F(0) = (0 * 4) + (1 * 3) + (2 * 2) + (3 * 6) = 0 + 3 + 4 + 18 = 25
# F(1) = (0 * 6) + (1 * 4) + (2 * 3) + (3 * 2) = 0 + 4 + 6 + 6 = 16
# F(2) = (0 * 2) + (1 * 6) + (2 * 4) + (3 * 3) = 0 + 6 + 8 + 9 = 23
# F(3) = (0 * 3) + (1 * 2) + (2 * 6) + (3 * 4) = 0 + 2 + 12 + 12 = 26

# 所以 F(0), F(1), F(2), F(3) 中的最大值是 F(3) = 26 。
