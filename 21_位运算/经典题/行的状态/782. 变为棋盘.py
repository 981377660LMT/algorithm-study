from collections import Counter
from itertools import chain
from typing import List, Sequence, Tuple


# 每次移动，你能任意交换两列或是两行的位置。
# !思路：要么等于第一行 要么和第一行完全相反
# board 是方阵，且行列数的范围是[2, 30]。


class Solution:
    def movesToChessboard(self, board: List[List[int]]) -> int:
        N = len(board)
        res = 0
        for counter in [Counter(map(tuple, board)), Counter(zip(*board))]:
            # 必须要两种状态 且至多差1
            if len(counter) != 2 or sorted(counter.values()) != [N // 2, (N + 1) // 2]:
                return -1
            # 状态必须交替间隔
            line1, line2 = counter.keys()
            if not all(a ^ b for a, b in zip(line1, line2)):
                return -1

            # 变为 0/1 开头
            starts = [int(line1.count(1) * 2 > N)] if N & 1 else [0, 1]
            res += (
                min(sum((i - num) & 1 for i, num in enumerate(line1, start=s)) for s in starts) // 2
            )

        return res


print(Solution().movesToChessboard(board=[[0, 1, 1, 0], [0, 1, 1, 0], [1, 0, 0, 1], [1, 0, 0, 1]]))
print(Solution().movesToChessboard(board=[[1, 1, 0], [0, 0, 1], [0, 0, 1]]))
# 输出: 2

# 解释:
# 一种可行的变换方式如下，从左到右：

# 0110     1010     1010
# 0110 --> 1010 --> 0101
# 1001     0101     1010
# 1001     0101     0101

# 第一次移动交换了第一列和第二列。
# 第二次移动交换了第二行和第三行。

