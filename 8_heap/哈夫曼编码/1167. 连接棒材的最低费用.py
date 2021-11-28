import heapq
from typing import List


class Solution:
    def connectSticks(self, sticks: List[int]) -> int:
        # 最小堆  哈夫曼树的思想，就是求最小cost
        minHeap = sticks
        heapq.heapify(minHeap)

        res = 0
        while len(minHeap) >= 2:
            x = heapq.heappop(minHeap)
            y = heapq.heappop(minHeap)
            res += x + y
            heapq.heappush(minHeap, x + y)
        return res


print(Solution().connectSticks([1, 8, 3, 5]))
# 输入：sticks = [1,8,3,5]
# 输出：30
# 解释：从 sticks = [1,8,3,5] 开始。
# 1. 连接 1 和 3 ，费用为 1 + 3 = 4 。现在 sticks = [4,8,5]
# 2. 连接 4 和 5 ，费用为 4 + 5 = 9 。现在 sticks = [9,8]
# 3. 连接 9 和 8 ，费用为 9 + 8 = 17 。现在 sticks = [17]
# 所有棒材已经连成一根，总费用 4 + 9 + 17 = 30

