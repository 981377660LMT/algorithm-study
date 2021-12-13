from typing import List
from functools import lru_cache
from itertools import accumulate

INF = 0x7FFFFFFF
# 看数据，不可能dfs
# 1 <= fruits.length <= 105
# 0 <= k <= 2 * 105


# 前缀和记录每个位置之前的所有水果数量,我们遍历所有位置,求出其最大的覆盖线段长度,并记录其最大的水果数量即可
class Solution:
    def maxTotalFruits(self, fruits: List[List[int]], startPos: int, k: int) -> int:
        # rightBound = max(i for i, _ in fruits)
        nums = [0] * (int(2e5 + 2))
        for i, c in fruits:
            nums[i] = c

        preSum = [0] + list(accumulate(nums))
        res = 0

        # print(preSum)
        # 贪心思路，结果最多只能掉头一次;枚举掉头的点 然后算出最大和的子数组
        for turn in range(len(nums)):
            # 从起点到掉头处的距离,从掉头处到重点的距离
            distBeforeTurn = abs(turn - startPos)
            distAfterTurn = k - distBeforeTurn
            if distBeforeTurn > k:
                continue

            # 左掉头/右掉头
            if turn <= startPos:
                end = turn + distAfterTurn
                rightMost = min(len(preSum) - 2, max(startPos, end))
                res = max(res, preSum[rightMost + 1] - preSum[turn])
            else:
                end = turn - distAfterTurn
                leftMost = max(0, min(startPos, end))
                res = max(res, preSum[turn + 1] - preSum[leftMost])

        return res


print(
    Solution().maxTotalFruits(
        fruits=[[0, 9], [4, 1], [5, 7], [6, 2], [7, 4], [10, 9]], startPos=5, k=4
    )
)
# 输出：14
# 解释：
# 可以移动最多 k = 4 步，所以无法到达位置 0 和位置 10 。
# 最佳路线为：
# - 在初始位置 5 ，摘到 7 个水果
# - 向左移动到位置 4 ，摘到 1 个水果
# - 向右移动到位置 6 ，摘到 2 个水果
# - 向右移动到位置 7 ，摘到 4 个水果
# 移动 1 + 3 = 4 步，共摘到 7 + 1 + 2 + 4 = 14 个水果
print(Solution().maxTotalFruits(fruits=[[2, 8], [6, 3], [8, 6]], startPos=5, k=4))
# 输出：9
# 解释：
# 最佳路线为：
# - 向右移动到位置 6 ，摘到 3 个水果
# - 向右移动到位置 8 ，摘到 6 个水果
# 移动 3 步，共摘到 3 + 6 = 9 个水果

print(Solution().maxTotalFruits(fruits=[[0, 3], [6, 4], [8, 5]], startPos=3, k=2))
print(Solution().maxTotalFruits(fruits=[[200000, 10000]], startPos=200000, k=200000))
