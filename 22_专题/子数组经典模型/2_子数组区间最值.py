from typing import List

INF = 0x3F3F3F3F


class Solution:
    def subArrayRanges(self, nums: List[int]) -> None:
        n = len(nums)
        minDp = [[INF] * n for _ in range(n)]
        maxDp = [[-INF] * n for _ in range(n)]

        for i in range(n):
            minDp[i][i] = nums[i]
            maxDp[i][i] = nums[i]
            for j in range(i + 1, n):
                minDp[i][j] = min(minDp[i][j - 1], nums[j])
                maxDp[i][j] = max(maxDp[i][j - 1], nums[j])

        print(maxDp[2][4])
        print(minDp[2][4])


print(Solution().subArrayRanges([1, 2, 3, 4, 5]))

