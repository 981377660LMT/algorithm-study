# 选择任意一个 nums[i] ，删除它并获得 nums[i] 的点数。
# !之后，你必须删除 所有 等于 nums[i] - 1 和 nums[i] + 1 的元素。
# 开始你拥有 0 个点数。返回你能通过这些操作获得的最大点数。
# 1 <= nums.length <= 2e4
# 1 <= nums[i] <= 1e4


# !思考:如果nums[i]<=1e9怎么做?


from typing import List
from collections import Counter


# 数轴上有n个位置,每个位置有分数scores[i]
# 不能偷相邻的位置(x-1,x+1),求最大分数
# !positions[i]<=1e9 且单调递增
def rob2(positions: List[int], scores: List[int]) -> int:
    if not positions:
        return 0

    n = len(positions)
    dp0, dp1 = 0, scores[0]
    for i in range(1, n):
        dist = positions[i] - positions[i - 1]
        if dist <= 1:
            dp0, dp1 = max(dp0, dp1), max(dp0 + scores[i], dp1)
        else:
            dp0, dp1 = max(dp0, dp1), max(dp0, dp1) + scores[i]
    return max(dp0, dp1)


class Solution:
    def deleteAndEarn(self, nums: List[int]) -> int:
        counter = Counter(nums)
        positions = sorted(counter)
        scores = [counter[v] * v for v in positions]
        return rob2(positions, scores)


print(Solution().deleteAndEarn([2, 2, 3, 3, 3, 4]))
