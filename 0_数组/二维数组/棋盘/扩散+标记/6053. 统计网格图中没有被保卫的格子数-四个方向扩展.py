from typing import List

# 请你返回空格子中，有多少个格子是 没被保卫 的。
# 前缀后缀统计(蜡烛间的盘子)
# 这种方法很差，可行但是代码量很大

DIRS4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]

# 从空格出发：前后缀解法，对每个空格，看上下左右离它最近的首卫是不是比墙近(非常繁琐)
# 从守卫出发：直接四个方向前进，碰到墙壁或者守卫就停下。这样为什么不会超时呢，因为每个点最多被每个行/列的守卫扫到两次，所以时间复杂度是 O(m*n) 的(更好的解法)
# 这种扩散+染色visited的解法在


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

        return sum(1 for r in range(row) for c in range(col) if matrix[r][c] == 0)

    def countUnguarded2(
        self, m: int, n: int, guards: List[List[int]], walls: List[List[int]]
    ) -> int:
        def check(r: int, c: int) -> bool:
            """(r,c)被保卫"""
            up1, right1, down1, left1 = (
                upGood[c][r],
                rightGood[r][c],
                downGood[c][r],
                leftGood[r][c],
            )

            up2, right2, down2, left2 = (
                upBad[c][r],
                rightBad[r][c],
                downBad[c][r],
                leftBad[r][c],
            )

            if up1 >= up2 and up1 != -1:
                return True
            if left1 >= left2 and left1 != -1:
                return True
            if right1 <= right2 and right1 != n:
                return True
            if down1 <= down2 and down1 != m:
                return True

            return False

        good = set(tuple(x) for x in guards)
        bad = set(tuple(x) for x in walls)

        leftGood = [[-1] * n for _ in range(m)]
        rightGood = [[n] * n for _ in range(m)]
        upGood = [[-1] * m for _ in range(n)]
        downGood = [[m] * m for _ in range(n)]

        leftBad = [[-1] * n for _ in range(m)]
        rightBad = [[n] * n for _ in range(m)]
        upBad = [[-1] * m for _ in range(n)]
        downBad = [[m] * m for _ in range(n)]

        for r in range(m):
            for c in range(n):
                if (r, c) in good:
                    leftGood[r][c] = c
                    upGood[c][r] = r
                elif (r, c) in bad:
                    leftBad[r][c] = c
                    upBad[c][r] = r
                if c - 1 >= 0:
                    leftGood[r][c] = max(leftGood[r][c], leftGood[r][c - 1])
                    leftBad[r][c] = max(leftBad[r][c], leftBad[r][c - 1])
                if r - 1 >= 0:
                    upGood[c][r] = max(upGood[c][r], upGood[c][r - 1])
                    upBad[c][r] = max(upBad[c][r], upBad[c][r - 1])

        for r in range(m - 1, -1, -1):
            for c in range(n - 1, -1, -1):
                if (r, c) in good:
                    rightGood[r][c] = c
                    downGood[c][r] = r
                elif (r, c) in bad:
                    rightBad[r][c] = c
                    downBad[c][r] = r
                if c + 1 < n:
                    rightGood[r][c] = min(rightGood[r][c], rightGood[r][c + 1])
                    rightBad[r][c] = min(rightBad[r][c], rightBad[r][c + 1])
                if r + 1 < m:
                    downGood[c][r] = min(downGood[c][r], downGood[c][r + 1])
                    downBad[c][r] = min(downBad[c][r], downBad[c][r + 1])

        res = 0
        for r in range(m):
            for c in range(n):
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

