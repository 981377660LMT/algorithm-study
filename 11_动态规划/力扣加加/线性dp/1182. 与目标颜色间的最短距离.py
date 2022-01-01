from typing import List

# 1 <= colors.length <= 5*10^4
# 1 <= colors[i] <= 3
# 查找从索引 i 到具有目标颜色 c 的元素之间的最短距离。
# dp[i][0]，dp[i][1]，dp[i][2]分别表示索引i到某种颜色的最短距离
# 用两遍扫描确定在i时的最短距离
class Solution:
    def shortestDistanceColor(self, colors: List[int], queries: List[List[int]]) -> List[int]:
        n = len(colors)
        dp = [[0x7FFFFFFF] * 3 for _ in range(n)]
        pre1 = pre2 = pre3 = -1

        for i in range(n):
            if colors[i] == 1:
                pre1 = i
            if colors[i] == 2:
                pre2 = i
            if colors[i] == 3:
                pre3 = i
            if pre1 != -1:
                dp[i][0] = min(dp[i][0], i - pre1)
            if pre2 != -1:
                dp[i][1] = min(dp[i][1], i - pre2)
            if pre3 != -1:
                dp[i][2] = min(dp[i][2], i - pre3)

        pre1 = pre2 = pre3 = -1
        for i in range(n - 1, -1, -1):
            if colors[i] == 1:
                pre1 = i
            if colors[i] == 2:
                pre2 = i
            if colors[i] == 3:
                pre3 = i
            if pre1 != -1:
                dp[i][0] = min(dp[i][0], pre1 - i)
            if pre2 != -1:
                dp[i][1] = min(dp[i][1], pre2 - i)
            if pre3 != -1:
                dp[i][2] = min(dp[i][2], pre3 - i)

        res = []
        for i, j in queries:
            if dp[i][j - 1] != 0x7FFFFFFF:
                res.append(dp[i][j - 1])
            else:
                res.append(-1)

        return res


print(
    Solution().shortestDistanceColor(
        colors=[1, 1, 2, 1, 3, 2, 2, 3, 3], queries=[[1, 3], [2, 2], [6, 1]]
    )
)
# 输出：[3,0,3]
# 解释：
# 距离索引 1 最近的颜色 3 位于索引 4（距离为 3）。
# 距离索引 2 最近的颜色 2 就是它自己（距离为 0）。
# 距离索引 6 最近的颜色 1 位于索引 3（距离为 3）。
