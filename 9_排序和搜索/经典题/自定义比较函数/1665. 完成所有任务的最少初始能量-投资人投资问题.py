from typing import List

# 给你一个任务数组 tasks ，其中 tasks[i] = [actuali, minimumi]
# actuali 是完成第 i 个任务 需要耗费 的实际能量。
# minimumi 是开始第 i 个任务前需要达到的最低能量。
# 1 <= tasks.length <= 105

# 你可以按照 任意顺序 完成任务。
# 请你返回完成所有任务的 最少 初始能量。


# 先易后难
# https://leetcode-cn.com/problems/minimum-initial-energy-to-finish-tasks/solution/chun-li-zi-jiang-jie-bang-zhu-li-jie-lao-9fgq/
class Solution:
    def minimumEffort(self, tasks: List[List[int]]) -> int:
        tasks = sorted(tasks, key=lambda t: t[1] - t[0])
        res = 0
        for actual, need in tasks:
            res = max(res + actual, need)
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
