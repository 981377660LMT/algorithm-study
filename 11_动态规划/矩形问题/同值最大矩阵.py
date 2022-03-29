class Solution:
    def solve(self, m):
        dp = [[1] * row for row in m]

        for i, row in enumerate(m[1:], start=1):
            for j, row in enumerate(row[1:], start=1):
                if m[i][j] == m[i - 1][j - 1] == m[i - 1][j] == m[i][j - 1]:
                    dp[i][j] = 1 + min(dp[i - 1][j - 1], dp[i - 1][j], dp[i][j - 1])

        return max(max(row) for row in dp)


# matrix = [
#     [1, 1, 3],
#     [1, 1, 3],
#     [5, 5, 5]
# ]
