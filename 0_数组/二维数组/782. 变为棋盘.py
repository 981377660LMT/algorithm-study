from typing import List

INF = 0x3F3F3F3F
# 思路：所有行全部参考第一行，所有列全部参考第一列
# board 是方阵，且行列数的范围是[2, 30]。
class Solution:
    def movesToChessboard(self, board: List[List[int]]) -> int:
        n = len(board)

        def cal(vals):
            """Return min moves to transform to chessboard."""
            first = vals[0]
            even = odd = 0
            for i, state in enumerate(vals):
                # 与第0行相同格子
                if first == state:
                    # 是第奇数行
                    if i & 1:
                        odd += 1
                    else:
                        even += 1
                # 与第一行不同格子:必须是与第一行全部相反
                elif first ^ state != (1 << n) - 1:
                    return INF

            res = INF
            # 2/4  3/5的情况 第一行不换 调整奇数行
            if n <= 2 * (even + odd) <= n + 1:
                res = min(res, odd)
            # 1/3 2/5的情况 => 第一行要换掉 调整偶数行
            if n - 1 <= 2 * (even + odd) <= n:
                res = min(res, even)
            return res

        rows, cols = [0] * n, [0] * n
        for i in range(n):
            for j in range(n):
                if board[i][j]:
                    rows[i] ^= 1 << j
                    cols[j] ^= 1 << i
        # print(rows, cols)
        res = cal(rows) + cal(cols)
        return res if res < INF else -1


print(Solution().movesToChessboard(board=[[0, 1, 1, 0], [0, 1, 1, 0], [1, 0, 0, 1], [1, 0, 0, 1]]))
# 输出: 2
# 解释:
# 一种可行的变换方式如下，从左到右：

# 0110     1010     1010
# 0110 --> 1010 --> 0101
# 1001     0101     1010
# 1001     0101     0101

# 第一次移动交换了第一列和第二列。
# 第二次移动交换了第二行和第三行。

