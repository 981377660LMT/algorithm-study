from typing import List

# 菱形和 指的是 grid 中一个正菱形 边界 上的元素之和。本题中的菱形必须为正方形旋转45度，且四个角都在一个格子当中。
# 请你按照 降序 返回 grid 中三个最大的 互不相同的菱形和.如果不同的和少于三个，则将它们全部返回。
# 枚举中心点，枚举边长
class Solution:
    def getBiggestThree(self, grid: List[List[int]]) -> List[int]:
        Row, Col = len(grid), len(grid[0])
        maxLen = min(Row, Col)
        val_set = set()

        for r in range(Row):
            for c in range(Col):
                val_set.add(grid[r][c])

        for r in range(Row):
            for c in range(Col):
                for Len in range(1, maxLen + 1):
                    U = r - Len
                    D = r + Len
                    L = c - Len
                    R = c + Len
                    if (0 <= U) and (D < Row) and (0 <= L) and (R < Col):
                        cur_sum = grid[r][L]  # 最左点
                        r0, c0 = r - 1, L + 1
                        while r0 > U:
                            cur_sum += grid[r0][c0]
                            r0 -= 1
                            c0 += 1

                        cur_sum += grid[U][c]  # 最上点
                        r0, c0 = U + 1, c + 1
                        while r0 < r:
                            cur_sum += grid[r0][c0]
                            r0 += 1
                            c0 += 1

                        cur_sum += grid[r][R]  # 最右点
                        r0, c0 = r + 1, R - 1
                        while r0 < D:
                            cur_sum += grid[r0][c0]
                            r0 += 1
                            c0 -= 1

                        cur_sum += grid[D][c]  # 最下点
                        r0, c0 = D - 1, c - 1
                        while r0 > r:
                            cur_sum += grid[r0][c0]
                            r0 -= 1
                            c0 -= 1

                        val_set.add(cur_sum)

                    else:
                        break

        a = list(val_set)
        a.sort(reverse=True)
        return a[:3]


print(
    Solution().getBiggestThree(
        grid=[
            [3, 4, 5, 1, 3],
            [3, 3, 4, 2, 3],
            [20, 30, 200, 40, 10],
            [1, 5, 5, 4, 1],
            [4, 3, 2, 2, 5],
        ]
    )
)
