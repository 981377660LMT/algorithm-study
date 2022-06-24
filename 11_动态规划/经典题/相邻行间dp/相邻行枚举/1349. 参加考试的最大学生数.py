from collections import defaultdict
from typing import List


# 如果座位是坏的（不可用），就用 '#' 表示；否则，用 '.' 表示。
# 学生可以看到左侧、右侧、左上、右上这四个方向上紧邻他的学生的答卷,但是看不到直接坐在他前面或者后面的学生的答卷
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
        def check(pre: int, cur: int) -> bool:
            """不能看到左侧、右侧、左上、右上这四个方向上紧邻他的学生的答卷"""
            s1, s2 = cur << 1, cur >> 1
            if cur & s1 or cur & s2 or pre & s1 or pre & s2:
                return False
            return True

        ROW, COL = len(seats), len(seats[0])
        gridStates = [0] * ROW
        for r in range(ROW):
            for c in range(COL):
                if seats[r][~c] == '#':
                    gridStates[r] |= 1 << c

        dp = defaultdict(int)
        for s in range(1 << COL):
            if s & gridStates[0]:
                continue
            if check(0, s):
                dp[s] = s.bit_count()

        for r in range(1, ROW):
            ndp = defaultdict(int)
            for pre in dp:
                for cur in range(1 << COL):
                    if cur & gridStates[r]:
                        continue
                    if check(pre, cur):
                        ndp[cur] = max(ndp[cur], dp[pre] + cur.bit_count())
            dp = ndp

        return max(dp.values())


print(
    Solution().maxStudents(
        seats=[
            ["#", ".", "#", "#", ".", "#"],
            [".", "#", "#", "#", "#", "."],
            ["#", ".", "#", "#", ".", "#"],
        ]
    )
)
