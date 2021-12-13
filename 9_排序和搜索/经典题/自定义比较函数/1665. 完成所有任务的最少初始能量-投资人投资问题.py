from typing import List

# 给你一个任务数组 tasks ，其中 tasks[i] = [actuali, minimumi]
# actuali 是完成第 i 个任务 需要耗费 的实际能量。
# minimumi 是开始第 i 个任务前需要达到的最低能量。
# 1 <= tasks.length <= 105

# 你可以按照 任意顺序 完成任务。
# 请你返回完成所有任务的 最少 初始能量。

# 即:有n个投资人,投资actual的能量,他认为你的公司值mimumal
# We have n investor,
# each investor can invest actual energy,
# and he think your company have minimum value.
# If you already have at least minimum - actual energy,
# the investor will invest you actual energy.
# Assuming you get all investment,
# how many energy you have at least in the end?

# => Need to sort by minimum - actual
class Solution:
    def minimumEffort(self, tasks: List[List[int]]) -> int:
        tasks = sorted(tasks, key=lambda t: t[1] - t[0])
        res = 0
        for a, m in tasks:
            res = max(res + a, m)
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
