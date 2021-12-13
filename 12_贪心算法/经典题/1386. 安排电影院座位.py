from typing import List
from collections import defaultdict

# 1 <= n <= 10^9
# 1 <= reservedSeats.length <= min(10*n, 10^4)
# 电影院的观影厅中有 n 行座位，行编号从 1 到 n ，且每一行内总共有 10 个座位，列编号从 1 到 10 。
# 请你返回 最多能安排多少个 4 人家庭
# 10 => 暗示状态压缩
class Solution:
    def maxNumberOfFamilies(self, n: int, reservedSeats: List[List[int]]) -> int:
        s1, s2, s3 = 0b00001111, 0b11000011, 0b11110000
        rowState = defaultdict(int)
        for row, col in reservedSeats:
            if 2 <= col <= 9:
                rowState[row] |= 1 << (col - 2)
        res = (n - len(rowState)) * 2
        for state in rowState.values():
            if any((state | s) == s for s in (s1, s2, s3)):
                res += 1
        return res


print(
    Solution().maxNumberOfFamilies(
        n=3, reservedSeats=[[1, 2], [1, 3], [1, 8], [2, 6], [3, 1], [3, 10]]
    )
)
# 输出：4
# 解释：上图所示是最优的安排方案，
# 总共可以安排 4 个家庭。蓝色的叉表示被预约的座位，橙色的连续座位表示一个 4 人家庭。
