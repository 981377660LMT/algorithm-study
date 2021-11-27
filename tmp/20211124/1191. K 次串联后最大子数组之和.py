from typing import List

# 1 <= arr.length <= 10^5
# 1 <= k <= 10^5
# 子数组长度可以是 0，在这种情况下它的总和也是 0。
class Solution:
    def kConcatenationMaxSum(self, arr: List[int], k: int) -> int:
        def maxSubArray(arr: List[int]) -> int:
            if not arr:
                return 0
            curMax = resMax = arr[0]
            for i in range(1, len(arr)):
                curMax = max(arr[i], curMax + arr[i])
                resMax = max(resMax, curMax)

            return max(resMax, 0)

        MOD = 10 ** 9 + 7
        arrSum = sum(arr)
        if k == 1:
            return maxSubArray(arr) % MOD
        if k == 2:
            return maxSubArray(arr * 2) % MOD
        if k > 2 and arrSum > 0:
            return ((k - 2) * arrSum + maxSubArray(arr * 2)) % MOD
        else:
            return maxSubArray(arr * 2) % MOD


print(Solution().kConcatenationMaxSum([1, -2, 1], 5))

