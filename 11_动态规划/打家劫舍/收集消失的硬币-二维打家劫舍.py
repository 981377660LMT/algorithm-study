# 每收集一个硬币，上下左右的硬币就会消失
# 求最大的分

# 二维打家劫舍
class Solution:
    def solve(self, matrix):
        def solve1D(row):
            """一维打家劫舍"""
            if len(row) == 1:
                return row[0]

            dp = [0] * len(row)
            dp[0] = row[0]
            dp[1] = max(row[0], row[1])
            for i in range(2, len(row)):
                dp[i] = max(dp[i - 1], dp[i - 2] + row[i])

            return dp[-1]

        if not any(matrix):
            return 0

        return solve1D([solve1D(row) for row in matrix])


print(Solution().solve(matrix=[[1, 7, 6, 5], [9, 9, 3, 1], [4, 8, 1, 2]]))
