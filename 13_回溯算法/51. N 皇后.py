from typing import List


class Solution:
    def solveNQueens(self, n: int) -> List[List[str]]:
        """
        返回所有不同的 n 皇后问题 的解，每个解是一个 n×n 的棋盘布局。
        使用三个位掩码分别记录已被占用的列、主对角线、反对角线位置，
        并维护一个长度为 n 的数组 `pos`，pos[row] = 列号。
        """

        def backtrack(row: int, columns: int, diagonals: int, anti_diagonals: int):
            # row: 当前处理的行号
            # columns: 列占用位图
            # diagonals: 主对角线占用位图（从上到下，向左移）
            # anti_diagonals: 反对角线占用位图（从上到下，向右移）
            if row == n:
                board = []
                for r in range(n):
                    line = ["."] * n
                    line[pos[r]] = "Q"
                    board.append("".join(line))
                res.append(board)
                return

            # 计算本行可放置的位置掩码（1 表示可放）
            available = (~(columns | diagonals | anti_diagonals)) & mask
            while available:
                bit = available & -available  # 取最低位的 1
                available &= available - 1  # 清除最低位的 1
                col = bit.bit_length() - 1  # 计算列号
                pos[row] = col
                # 递归下一行，diagonals 左移 1，anti_diagonals 右移 1
                backtrack(
                    row + 1, columns | bit, (diagonals | bit) << 1, (anti_diagonals | bit) >> 1
                )

        mask = (1 << n) - 1
        pos = [0] * n
        res: List[List[str]] = []
        backtrack(0, 0, 0, 0)
        return res


if __name__ == "__main__":
    sol = Solution()
    for n in [4, 5]:
        print(f"n = {n} 解共有 {len(sol.solveNQueens(n))} 种：")
        for board in sol.solveNQueens(n):
            for line in board:
                print(line)
            print()
