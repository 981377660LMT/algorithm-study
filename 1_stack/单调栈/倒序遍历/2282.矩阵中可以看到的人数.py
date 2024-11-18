# 2282.矩阵中可以看到的人数
# 每个人可以看到他右方/下方的人(如果中间的人比他们都矮)

from typing import List


def canSeePersonsCount(heights: List[int]) -> List[int]:
    res = [0] * len(heights)
    stack = []
    for i in range(len(heights) - 1, -1, -1):
        while stack and stack[-1] < heights[i]:  # 必须要严格小于才不会被挡住
            res[i] += 1
            stack.pop()
        if stack:
            res[i] += 1
        while stack and stack[-1] == heights[i]:  # !相等的全部出栈，只能算一次
            stack.pop()
        stack.append(heights[i])
    return res


class Solution:
    def seePeople(self, matrix: List[List[int]]) -> List[List[int]]:
        m, n = len(matrix), len(matrix[0])
        res = [[0] * n for _ in range(m)]
        for r, row in enumerate(matrix):
            cur = canSeePersonsCount(row)
            for c, num in enumerate(cur):
                res[r][c] += num

        for c, col in enumerate(zip(*matrix)):
            cur = canSeePersonsCount(list(col))
            for r, num in enumerate(cur):
                res[r][c] += num

        return res


print(Solution().seePeople([[4, 2, 1, 1, 3]]))

# [[2,2,1,1,0]]
