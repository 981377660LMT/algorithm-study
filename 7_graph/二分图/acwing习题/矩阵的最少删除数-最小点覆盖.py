# 有一个01对称矩阵，你可以选择每次消除第i行第i列的所有1（不是一个1)，是整行整列都消除，
# 问最少需要几次操作才能把所有的1消除

# 求解二分图的最小点覆盖
# 一边是行坐标，一边是列坐标
# 然后二分图的边就是要覆盖(删)的点
# 选择最少的点数可以覆盖所有边
# 每一次删除相当于进行了两次点覆盖操作
# 二分图最小点覆盖可以直接转成求解二分图最大匹配来求


from typing import List
from hungarian import Hungarian


def minDelete(grid: List[List[int]]) -> int:
    """求二分图的最小点覆盖"""
    ROW, COL = len(grid), len(grid[0])
    H = Hungarian(ROW * COL, ROW * COL)
    for r in range(ROW):
        for c in range(COL):
            if grid[r][c] == 1:
                H.addEdge(r, c)

    return H.work()


print(minDelete(grid=[[1, 0, 1, 1], [0, 0, 1, 1], [1, 1, 0, 0], [1, 1, 0, 1]]))


# TODO
