# 有一个01对称矩阵，你可以选择每次消除第i行第i列的所有1（不是一个1)，是整行整列都消除，
# 问最少需要几次操作才能把所有的1消除

# 求解二分图的最小点覆盖
# 一边是行坐标，一边是列坐标
# 然后二分图的边就是要覆盖(删)的点
# 选择最少的点数可以覆盖所有边
# 每一次删除相当于进行了两次点覆盖操作
# 二分图最小点覆盖可以直接转成求解二分图最大匹配来求


from pprint import pprint
from typing import List
from collections import defaultdict
from hungarian import hungarian


def minDelete(grid: List[List[int]]) -> int:
    adjMap = defaultdict(set)
    res = 0
    remove = set()
    for i in range(len(grid)):
        if grid[i][i] == 1:
            # 必须要删的位置
            res += 1
            remove.add(i)

    for i in range(len(grid)):
        if i in remove:
            continue
        for j in range(len(grid[0])):
            if j in remove:
                continue
            if grid[i][j] == 1:
                adjMap[i].add(int(1e9) + j)
                adjMap[int(1e9) + j].add(i)
    pprint(grid, width=20)
    print(adjMap)
    return hungarian(adjMap) + res


print(minDelete(grid=[[1, 0, 1, 1], [0, 0, 1, 1], [1, 1, 0, 0], [1, 1, 0, 1]]))

