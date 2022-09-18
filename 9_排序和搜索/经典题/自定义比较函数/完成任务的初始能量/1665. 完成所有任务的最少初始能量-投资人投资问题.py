"""1665. 完成所有任务的最少初始能量-投资人投资问题

给你一个任务数组 tasks ，其中 tasks[i] = [actuali, costi]
actuali 是完成第 i 个任务 需要耗费 的实际能量。
costi 是开始第 i 个任务前需要达到的最低能量。
1 <= tasks.length <= 1e5 , 1 <= actuali <= minimumi <= 104


你可以按照 任意顺序 完成任务。
请你返回完成所有任务的 `最少` 初始能量。
"""

from typing import List


# !前面cashback越大就越容易
class Solution:
    def minimumEffort(self, tasks: List[List[int]]) -> int:
        bad, good = [], []
        for auctual, cost in tasks:
            cashback = cost - auctual
            if auctual > 0:
                bad.append((cost, cashback))
            else:
                good.append((cost, cashback))

        bad.sort(key=lambda x: x[1], reverse=True)  # !前面cashback越大就越容易
        good.sort(key=lambda x: x[0])  # !前面cost越小就越容易
        nums = bad + good

        res, cur = 0, 0
        for cost, cashback in nums:
            if cur < cost:
                res += cost - cur
                cur = cost
            cur -= cost
            cur += cashback
        return res


print(Solution().minimumEffort(tasks=[[1, 3], [2, 4], [10, 11], [10, 12], [8, 9]]))
# 输出：32
# 解释：
# 一开始有 32 能量，我们按照如下顺序完成任务：
#     - 完成第 1 个任务，剩余能量为 32 - 1 = 31 。
#     - 完成第 2 个任务，剩余能量为 31 - 2 = 29 。
#     - 完成第 3 个任务，剩余能量为 29 - 10 = 19 。
#     - 完成第 4 个任务，剩余能量为 19 - 10 = 9 。
#     - 完成第 5 个任务，剩余能量为 9 - 8 = 1 。
