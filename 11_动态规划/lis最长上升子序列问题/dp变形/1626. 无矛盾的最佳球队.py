from typing import List

# 你想组合一支总体得分最高的球队。球队的得分是球队中所有球员的分数 总和 。
# 如果一名年龄较小球员的分数 严格大于 一名年龄较大的球员，则存在矛盾。同龄球员之间不会发生矛盾。
# 请你返回 所有可能的无矛盾球队中得分最高那支的分数 。

# 1 <= scores.length, ages.length <= 1000
class Solution:
    def bestTeamScore(self, scores: List[int], ages: List[int]) -> int:
        n = len(scores)
        # 对年龄排序后寻找LIS
        team = sorted(zip(ages, scores))

        dp = [score for _, score in team]
        for i in range(1, n):
            for j in range(i):
                if team[i][1] >= team[j][1]:
                    dp[i] = max(dp[i], dp[j] + team[i][1])
        return max(dp)


print(Solution().bestTeamScore(scores=[1, 2, 3, 5], ages=[8, 9, 10, 1]))
# 输出：6

