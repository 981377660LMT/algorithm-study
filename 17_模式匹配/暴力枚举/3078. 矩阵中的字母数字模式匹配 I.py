# 3078. 矩阵中的字母数字模式匹配 I
# https://leetcode.cn/problems/match-alphanumerical-pattern-in-matrix-i/

# 给定一个二维整数矩阵 board 和一个二维字符矩阵 pattern。其中 0 <= board[r][c] <= 9 并且 pattern 的每个元素是一个数字或一个小写英文字母。
# 你的任务是找到 匹配 board 的子矩阵 pattern。
# 如果我们能用一些数字（每个 不同 的字母对应 不同 的数字）替换 pattern 中包含的字母使得结果矩阵与整数矩阵 part 相同，我们称整数矩阵 part 与 pattern 匹配。换句话说，
# 这两个矩阵具有相同的维数。
# 如果 pattern[r][c] 是一个数字，那么 part[r][c] 必须是 相同的 数字。
# 如果 pattern[r][c] 是一个字母 x：
# 对于每个 pattern[i][j] == x，part[i][j] 一定与 part[r][c] 相同。
# 对于每个 pattern[i][j] != x，part[i][j] 一定与 part[r][c] 不同。
# 返回一个长度为 2 的数组，包含匹配 pattern 的 board 的子矩阵左上角的行号和列号。
# 如果有多个这样的子矩阵，返回行号更小的子矩阵。如果依然有多个，则返回列号更小的子矩阵。
# 如果没有符合的答案，返回 [-1, -1]。
# 1 <= board.length <= 50
# 1 <= board[i].length <= 50
# 0 <= board[i][j] <= 9
# 1 <= pattern.length <= 50
# 1 <= pattern[i].length <= 50
# pattern[i][j] 表示为一个数字的字符串或一个小写英文字母。

from typing import List


class Solution:
    def findPattern(self, board: List[List[int]], pattern: List[str]) -> List[int]:
        row, col = len(board), len(board[0])
        pRow, pCol = len(pattern), len(pattern[0])

        def check(r: int, c: int) -> bool:
            mp, revMp = dict(), dict()
            for i in range(r, r + pRow):
                for j in range(c, c + pCol):
                    v1, v2 = board[i][j], pattern[i - r][j - c]
                    if v2.isdigit():
                        if v1 != int(v2):
                            return False
                    else:
                        in1, in2 = v1 in mp, v2 in revMp
                        if in1 and in2:
                            if mp[v1] != v2 or revMp[v2] != v1:
                                return False
                        elif in1 and not in2:
                            return False
                        elif not in1 and in2:
                            return False
                        elif not in1 and not in2:
                            mp[v1], revMp[v2] = v2, v1
            return True

        for r in range(row - pRow + 1):
            for c in range(col - pCol + 1):
                if check(r, c):
                    return [r, c]
        return [-1, -1]
