# You want to pick a number from each row. For each 0 ≤ r < n - 1 the cost for picking matrix[r][j] and matrix[r + 1][k] is abs(k - j).
# 去除绝对值，合并相同的变量即可
# dp[i][j] = dp[i - 1][k] + (j - k) => (dp[i-1][k]-k)+j  # k<=j
# dp[i][j]=dp[i-1][k]+(k-j) => (dp[i-1][k]+k)-j # k>=j

INF = int(1e20)


class Solution:
    def solve(self, matrix):
        dp = matrix[0]
        for i in range(1, len(matrix)):
            ndp = [-INF] * len(matrix[i])
            left = -INF
            for j in range(len(matrix[i])):
                left = max(left, dp[j])
                ndp[j] = max(ndp[j], left + matrix[i][j])
                left -= 1
            right = -INF
            for j in range(len(matrix[i]) - 1, -1, -1):
                right = max(right, dp[j])
                ndp[j] = max(ndp[j], right + matrix[i][j])
                right -= 1
            dp = ndp
        return max(dp)
