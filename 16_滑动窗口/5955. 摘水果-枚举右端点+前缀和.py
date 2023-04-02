# 摘水果
# 从任何位置，你可以选择 向左或者向右 走。
# 在 x 轴上每移动 一个单位 ，就记作 一步 。
# 你总共可以走 最多 k 步。你每达到一个位置，都会摘掉全部的水果，水果也将从该位置消失（不会再生）。
# 返回你可以摘到水果的 最大总数 。
# 1 <= fruits.length <= 105
# 0 <= k <= 2 * 105

# !枚举最右端的点,看最左端可以到达哪个位置
# !调头一次的最短距离=min(2*leftMax+rightMax,2*rightMax+leftMax)


from collections import defaultdict
from typing import List


class Solution:
    def maxTotalFruits(self, fruits: List[List[int]], startPos: int, k: int) -> int:
        def calDist(start: int, left: int, right: int) -> int:
            """从start出发,遍历[leftMost,rightMost]区间的最短距离(最多调头一次)"""
            leftMax = max(0, start - left)
            rightMax = max(0, right - start)
            return min(2 * leftMax + rightMax, 2 * rightMax + leftMax)

        counter = defaultdict(int)  # 记录每个位置的水果数量
        for pos, count in fruits:
            counter[pos] += count
        keys = sorted(counter)

        n = len(keys)
        res, left, curSum = 0, 0, 0
        for right in range(n):
            curSum += counter[keys[right]]
            while left <= right and calDist(startPos, keys[left], keys[right]) > k:
                curSum -= counter[keys[left]]
                left += 1
            res = max(res, curSum)
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
