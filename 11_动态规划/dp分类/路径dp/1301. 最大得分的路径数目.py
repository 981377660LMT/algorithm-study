from typing import List

# 你从数组最右下方的字符 'S' 出发。
# 你的目标是到达数组最左上角的字符 'E'
# 在每一步移动中，你可以向上、向左或者左上方移动，可以移动的前提是到达的格子没有障碍。
# 请你返回一个列表，包含两个整数：第一个整数是 「得分」 的最大值，第二个整数是得到最大得分的方案数
# 2 <= board.length == board[i].length <= 100

# dp[x][y][0] is the maximum value to this cell,
# dp[x][y][1] is the number of paths.
MOD = int(1e9 + 7)
INF = 0x7FFFFFFF


class Solution:
    def pathsWithMaxScore(self, board: List[str]) -> List[int]:
        n = len(board)
        dp = [[[-INF, 0] for _ in range(n + 1)] for _ in range(n + 1)]
        dp[n - 1][n - 1] = [0, 1]

        for r in reversed(range(n)):
            for c in reversed(range(n)):
                if board[r][c] in 'XS':
                    continue
                # 全部从后往前
                for dr, dc in [[0, 1], [1, 0], [1, 1]]:
                    if dp[r][c][0] < dp[r + dr][c + dc][0]:
                        dp[r][c] = [dp[r + dr][c + dc][0], 0]
                    if dp[r][c][0] == dp[r + dr][c + dc][0]:
                        dp[r][c][1] += dp[r + dr][c + dc][1]
                dp[r][c][0] += int(board[r][c]) if (r, c) != (0, 0) else 0
        return [dp[0][0][0] if dp[0][0][1] else 0, dp[0][0][1] % MOD]


print(Solution().pathsWithMaxScore(["E23", "2X2", "12S"]))
# 输出：[7,1]
