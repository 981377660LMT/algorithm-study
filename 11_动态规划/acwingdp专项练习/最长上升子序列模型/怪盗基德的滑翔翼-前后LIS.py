# 当起点和方向确定了之后，问题就变成
# a[i]：求a[i]为终点的最长上升子序列
from typing import List

# n<=1000


class Solution:
    def Kiddo(self, buildings: List[int]):  # 可以用二分查找
        def getLIS(nums: List[int]) -> int:
            res = 0
            dp = [1 for _ in range(len(nums))]
            for i in range(len(dp)):
                for j in range(i):
                    if nums[i] > nums[j]:
                        dp[i] = max(dp[i], dp[j] + 1)
                        res = max(res, dp[i])
            return res

        return max(getLIS(buildings), getLIS(buildings[::-1]))


if __name__ == '__main__':
    res = []
    solution = Solution()
    group = int(input())
    for i in range(group):
        num = int(input())
        buildings = list(map(int, input().split()))
        res.append(solution.Kiddo(buildings))

    for j in res:
        print(j)
