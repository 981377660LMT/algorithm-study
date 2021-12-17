from typing import List


# 思路：从大到小排序，选择奇数的数量要是偶数个
class Solution:
    def largestEvenSum(self, nums: List[int], k: int) -> int:
        odd, even = [0], [0]
        # 奇偶前缀和分家
        for x in sorted(nums, reverse=True):
            if x & 1:
                odd.append(odd[-1] + x)
            else:
                even.append(even[-1] + x)

        best = -1
        # 枚举奇数个数
        n, m = len(odd), len(even)
        for x in range(0, n, 2):
            if 0 <= k - x < m:
                best = max(best, odd[x] + even[k - x])
        return best


print(Solution().largestEvenSum(nums=[4, 1, 5, 3, 1], k=3))
# Output: 12
# Explanation:
# The subsequence with the largest possible even sum is [4,5,3]. It has a sum of 4 + 5 + 3 = 12.

