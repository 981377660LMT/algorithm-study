class Solution:
    def totalNQueens(self, n: int) -> int:
        """
        返回 n 皇后不同解的数量。
        使用三个位掩码分别记录已被占用的列、主对角线、反对角线位置，
        每一行用 available 计算出可放置的位置，然后逐一递归。
        """

        def backtrack(row: int, columns: int, diagonals: int, anti_diagonals: int):
            nonlocal res
            if row == n:
                res += 1
                return

            available = (~(columns | diagonals | anti_diagonals)) & mask
            while available:
                position = available & -available
                available &= available - 1
                # 递归到下一行，更新三种“攻击”掩码
                # columns | position：占用该列
                # (diagonals | position) << 1：主对角线向下移一位
                # (anti_diagonals | position) >> 1：反对角线向下移一位
                backtrack(
                    row + 1,
                    columns | position,
                    (diagonals | position) << 1,
                    (anti_diagonals | position) >> 1,
                )

        mask = (1 << n) - 1
        res = 0
        backtrack(0, 0, 0, 0)
        return res


if __name__ == "__main__":
    for n in range(1, 11):
        print(f"n={n} 方案数 =", Solution().totalNQueens(n))
