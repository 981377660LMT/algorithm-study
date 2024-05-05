# 1953. 你可以工作的最大周数
# https://leetcode.cn/problems/maximum-number-of-weeks-for-which-you-can-work/description/
#
# 你可以按下面两个规则参与项目中的工作：
# - 每周，你将会完成 某一个 项目中的 恰好一个 阶段任务。你每周都 必须 工作。
# - 在 连续的 两周中，你 不能 参与并完成同一个项目中的两个阶段任务。
# 一旦所有项目中的全部阶段任务都完成，或者仅剩余一个阶段任务都会导致你违反上面的规则，那么你将 停止工作 。
# 返回在不违反上面规则的情况下你 最多 能工作多少周。

from typing import List


class Solution:
    def numberOfWeeks(self, milestones: List[int]) -> int:
        max_ = max(milestones)
        rest = sum(milestones) - max_
        if max_ > rest:
            return rest * 2 + 1
        else:
            return max_ + rest
