from collections import defaultdict
from typing import List

# 你想组合一支总体得分最高的球队。球队的得分是球队中所有球员的分数 总和 。
# 如果一名年龄较小球员的分数 严格大于 一名年龄较大的球员，则存在矛盾。同龄球员之间不会发生矛盾。
# 请你返回 所有可能的无矛盾球队中得分最高那支的分数 。

# 1 <= scores.length, ages.length <= 1000
# scores.length == ages.length
# 1 <= scores[i] <= 106
# 1 <= ages[i] <= 1000

# !weighted LIS


class Solution:
    def bestTeamScore1(self, scores: List[int], ages: List[int]) -> int:
        """O(nlogn) 树状数组记录"""
        dp = BIT3(1010)  # 每个年龄对应的最大分数
        team = sorted(zip(scores, ages))  # 按(分数,年龄)排序
        for score, age in team:
            preMax = dp.query(age)  # 前面年龄的最大分数
            dp.update(age, preMax + score)  # 更新当前年龄的最大分数
        return dp.query(1010)

    def bestTeamScore2(self, scores: List[int], ages: List[int]) -> int:
        """O(n^2)"""
        n = len(scores)
        # 对年龄排序后寻找LIS
        team = sorted(zip(ages, scores))

        dp = [score for _, score in team]
        for i in range(1, n):
            for j in range(i):
                if team[i][1] >= team[j][1]:
                    dp[i] = max(dp[i], dp[j] + team[i][1])
        return max(dp)


class BIT3:
    """单点修改 维护`前缀区间`最大值"""

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    def update(self, index: int, target: int) -> None:
        """将后缀区间`[index,size]`的最大值更新为target"""
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] = max(self.tree[index], target)
            index += index & -index

    def query(self, index: int) -> int:
        """查询前缀区间`[1,index]`的最大值"""
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res = max(res, self.tree[index])
            index -= index & -index
        return res


print(Solution().bestTeamScore1(scores=[1, 2, 3, 5], ages=[8, 9, 10, 1]))
print(Solution().bestTeamScore2(scores=[1, 2, 3, 5], ages=[8, 9, 10, 1]))
# 输出：6
# [1,1,1,1,1,1,1,1,1,1]
# [811,364,124,873,790,656,581,446,885,134]
print(
    Solution().bestTeamScore1(
        scores=[1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
        ages=[811, 364, 124, 873, 790, 656, 581, 446, 885, 134],
    )
)
