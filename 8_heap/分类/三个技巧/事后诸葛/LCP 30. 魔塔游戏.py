from typing import List
from heapq import heappush, heappop

# 请返回小扣最少需要调整几次，才能顺利访问所有房间
# 小扣初始血量为 1，且无上限

# 贪心
# !当不行了时候，把pq的最小的负数扔到最后
# 用堆来实现维护当前最小的负数


class Solution:
    def magicTower(self, nums: List[int]) -> int:
        if sum(nums) + 1 <= 0:
            return -1

        hp, res = 1, 0
        pq = []

        for num in nums:
            if num < 0:
                heappush(pq, num)
                if hp + num <= 0:
                    res += 1
                    hp -= heappop(pq)

            hp += num

        return res


print(Solution().magicTower([100, 100, 100, -250, -60, -140, -50, -50, 100, 150]))
# 输出：1
# 解释：初始血量为 1。至少需要将 nums[3] 调整至访问顺序末尾以满足要求。
