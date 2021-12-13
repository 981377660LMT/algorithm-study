from typing import List


class Solution:
    def numberOfWeeks(self, milestones: List[int]) -> int:
        maxVal = max(milestones)
        restSum = sum(milestones) - maxVal
        if maxVal > restSum:
            return restSum * 2 + 1
        else:
            return maxVal + restSum


print(Solution().numberOfWeeks(milestones=[5, 2, 1]))
# 输出：7
# 解释：一种可能的情形是：
# - 第 1 周，你参与并完成项目 0 中的一个阶段任务。
# - 第 2 周，你参与并完成项目 1 中的一个阶段任务。
# - 第 3 周，你参与并完成项目 0 中的一个阶段任务。
# - 第 4 周，你参与并完成项目 1 中的一个阶段任务。
# - 第 5 周，你参与并完成项目 0 中的一个阶段任务。
# - 第 6 周，你参与并完成项目 2 中的一个阶段任务。
# - 第 7 周，你参与并完成项目 0 中的一个阶段任务。
# 总周数是 7 。
# 注意，你不能在第 8 周参与完成项目 0 中的最后一个阶段任务，因为这会违反规则。
# 因此，项目 0 中会有一个阶段任务维持未完成状态。
