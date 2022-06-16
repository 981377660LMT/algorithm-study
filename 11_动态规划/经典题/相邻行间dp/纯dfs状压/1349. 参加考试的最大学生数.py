from typing import List
from functools import lru_cache

# 如果座位是坏的（不可用）学生可以看到左侧、右侧、左上、右上这四个方向上紧邻他的学生的答卷，，就用 '#' 表示；否则，用 '.' 表示。
# 学生可以看到左侧、右侧、左上、右上这四个方向上紧邻他的学生的答卷，
# 请你计算并返回该考场可以容纳的一起参加考试且无法作弊的最大学生人数。

# 1 <= m <= 8
# 1 <= n <= 8

# 每一行只与上一行有关，同时每一行最多2^8  种状态，我们自然想到进行状态压缩DP。
# dp[row][state]=max(dp[row−1][last]+state.count())
# 本行的合法性：不能把学生安排在坏座位上；不能有相邻的学生
# 两行之间的合法性：如果第一行某个位置安排了学生，则下一行斜向的两个位置不能安排学生
# 最后的结果就是max(dp[m][state])


class Solution:
    def maxStudents(self, seats: List[List[str]]) -> int:
        @lru_cache(None)
        def dfs(rowIndex: int, preState: int) -> int:
            if rowIndex == ROW:
                return 0

            res = 0
            for curState in range(1 << COL):
                if not check(rowIndex, preState, curState):
                    continue
                res = max(res, curState.bit_count() + dfs(rowIndex + 1, curState))
            return res

        def check(row: int, preState: int, curState: int) -> bool:
            # 不能有相邻的1
            if curState & (curState << 1) or curState & (curState >> 1):
                return False
            # 不能看到左上、右上方
            if preState & (curState << 1) or preState & (curState >> 1):
                return False
            # 不能有坏座位
            for i, char in enumerate(seats[row]):
                if char == '#' and (curState >> i) & 1:
                    return False
            return True

        ROW, COL = len(seats), len(seats[0])
        return dfs(0, 0)


print(
    Solution().maxStudents(
        seats=[
            ["#", ".", "#", "#", ".", "#"],
            [".", "#", "#", "#", "#", "."],
            ["#", ".", "#", "#", ".", "#"],
        ]
    )
)
