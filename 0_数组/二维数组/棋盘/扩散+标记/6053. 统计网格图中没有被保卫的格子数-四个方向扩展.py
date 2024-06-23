# 请你返回空格子中，有多少个格子是 没被保卫 的。
# 前缀后缀统计(蜡烛间的盘子)
# 这种方法很差，可行但是代码量很大


# 从空格出发：前后缀解法，对每个空格，看上下左右离它最近的首卫是不是比墙近(非常繁琐)
# !从守卫出发：直接四个方向前进，碰到墙壁或者守卫就停下。这样为什么不会超时呢，因为每个点最多被每个行/列的守卫扫到两次，所以时间复杂度是 O(m*n) 的(更好的解法)
# 这种扩散+染色visited的解法在


from typing import List
from findNearestSpecial import findNearestSpecial

DIRS4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]


class Solution:
    def countUnguarded(
        self, row: int, col: int, guards: List[List[int]], walls: List[List[int]]
    ) -> int:
        matrix = [[0] * col for _ in range(row)]
        for r, c in guards:
            matrix[r][c] = 1
        for r, c in walls:
            matrix[r][c] = 2

        for r in range(row):
            for c in range(col):
                if matrix[r][c] != 1:
                    continue
                for dr, dc in DIRS4:
                    nr, nc = r + dr, c + dc
                    while 0 <= nr < row and 0 <= nc < col and matrix[nr][nc] not in (1, 2):
                        matrix[nr][nc] = 3
                        nr, nc = nr + dr, nc + dc

        return sum(matrix[r][c] == 0 for r in range(row) for c in range(col))

    def countUnguarded2(
        self, row: int, col: int, guards: List[List[int]], walls: List[List[int]]
    ) -> int:
        def check(r: int, c: int) -> bool:
            """(r,c)被保卫"""
            top1, bottom1, left1, right1 = nearestGood[r][c]
            top2, bottom2, left2, right2 = nearestBad[r][c]
            if top1 >= top2 and top1 != -1:
                return True
            if left1 >= left2 and left1 != -1:
                return True
            if right1 <= right2 and right1 != col:
                return True
            if bottom1 <= bottom2 and bottom1 != row:
                return True
            return False

        good = set(tuple(x) for x in guards)
        nearestGood = findNearestSpecial(row, col, lambda x: x in good)
        bad = set(tuple(x) for x in walls)
        nearestBad = findNearestSpecial(row, col, lambda x: x in bad)

        res = 0
        for r in range(row):
            for c in range(col):
                if (r, c) not in good and (r, c) not in bad:
                    if not check(r, c):
                        res += 1
        return res


print(
    Solution().countUnguarded(
        row=4, col=6, guards=[[0, 0], [1, 1], [2, 3]], walls=[[0, 1], [2, 2], [1, 4]]
    )
)
# print(Solution().countUnguarded(m=3, n=3, guards=[[1, 1]], walls=[[0, 1], [1, 0], [2, 1], [1, 2]]))
