from typing import List

# 1 <= moves.length <= 9
# moves 里没有重复的元素。
check = [
    0b111000000,
    0b000111000,
    0b000000111,
    0b100100100,
    0b010010010,
    0b001001001,
    0b100010001,
    0b001010100,
]


class Solution:
    def tictactoe(self, moves: List[List[int]]) -> str:
        a = sum(1 << (i * 3 + j) for i, j in moves[::2])
        b = sum(1 << (3 * i + j) for i, j in moves[1::2])
        for c in check:
            if a & c == c:
                return 'A'
            if b & c == c:
                return 'B'

        return ['Pending', 'Draw'][len(moves) == 9]


print(Solution().tictactoe(moves=[[0, 0], [2, 0], [1, 1], [2, 1], [2, 2]]))
# 解释："A" 获胜
