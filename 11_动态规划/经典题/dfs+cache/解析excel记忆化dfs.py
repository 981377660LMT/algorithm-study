from functools import lru_cache
from operator import add, sub, mul, truediv
from typing import List

# n, m ≤ 26
# 你可以假设excel中无环

OPTIONS = {'+': add, '-': sub, '*': mul, '/': truediv}


class Solution:
    def solve(self, matrix: List[List[str]]) -> List[List[str]]:
        def resolve(string: str) -> int:
            try:
                return int(string)
            except ValueError:
                return parse(int(string[1:]) - 1, ord(string[0]) - ord("A"))

        @lru_cache(None)
        def parse(r, c) -> int:
            try:
                return int(matrix[r][c])
            except ValueError:
                opt = ''
                exp = matrix[r][c]
                if exp[0] == '=':
                    for char in exp[2:]:
                        if char in OPTIONS:
                            opt = char
                            break
                    a, b = exp[1:].split(opt)
                    aRes, bRes = resolve(a), resolve(b)
                    return OPTIONS[opt](aRes, bRes)
                else:
                    return parse(int(exp[1:]) - 1, ord(exp[0]) - ord("A"))

        row, col = len(matrix), len(matrix[0])
        for r in range(row):
            for c in range(col):
                matrix[r][c] = str(parse(r, c))

        return matrix


print(
    Solution().solve(
        matrix=[["B1", "-2", "=6+0"], ["1", "=A3+A3", "=A2-A1"], ["=-2+B1", "=C1+10", "C2"]]
    )
)

# # [
# #     ["-2", "-2", "6"],
# #     ["1", "-8", "3"],
# #     ["-4", "16", "3"]
# # ]

# #     A     B    C
# # 1 [[B1,   2,    0],
# # 2  [3,    5,  =A2+A1]]
