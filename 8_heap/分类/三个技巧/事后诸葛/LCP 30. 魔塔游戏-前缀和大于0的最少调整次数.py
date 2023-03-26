# !LCP 30. 魔塔游戏-前缀和大于0的最少调整次数
# 请返回小扣最少需要调整几次，才能顺利访问所有房间
# !小扣初始血量为 1，且无上限

# 贪心
# !当不行了时候，把pq的最小的负数扔到最后
# 用堆来实现维护当前最小的负数


from typing import List
from heapq import heappush, heappop


def magicTower(nums: List[int]) -> int:
    if sum(nums) + 1 <= 0:
        return -1
    hp, res = 1, 0  # 小扣初始血量为 1，且无上限
    pq = []
    for num in nums:
        if num < 0:
            heappush(pq, num)
            if hp + num <= 0:  # !当不行了时候，把pq的最小的负数扔到最后
                res += 1
                hp -= heappop(pq)

        hp += num
    return res


class Solution:
    def makePrefSumNonNegative(self, nums: List[int]) -> int:
        # https://leetcode.cn/problems/make-the-prefix-sum-non-negative/
        # 初始和为0,求最小的调整次数使得前缀和始终大于等于0
        if sum(nums) < 0:
            return -1

        curSum, res = 0, 0
        pq = []
        for num in nums:
            if num < 0:
                heappush(pq, num)
                if curSum + num < 0:
                    res += 1
                    curSum -= heappop(pq)
            curSum += num

        return res


print(magicTower([100, 100, 100, -250, -60, -140, -50, -50, 100, 150]))
# 输出：1
# 解释：初始血量为 1。至少需要将 nums[3] 调整至访问顺序末尾以满足要求。
