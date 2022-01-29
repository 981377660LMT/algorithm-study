from typing import List
from collections import defaultdict, deque

# 每一次操作时，打印机会用同一种颜色打印一个矩形的形状，每次打印会覆盖矩形对应格子里原本的颜色。
# 一旦矩形根据上面的规则使用了一种颜色，那么` 相同的颜色不能再被使用 。`

# 给你一个初始没有颜色的 m x n 的矩形 targetGrid
# 如果你能按照上述规则打印出矩形 targetGrid ，请你返回 true ，否则返回 false 。

# 单向图不允许形成环


class Solution:
    def isPrintable(self, targetGrid: List[List[int]]) -> bool:
        row, col = len(targetGrid), len(targetGrid[0])
        # 遍历所有点，计算出每个值矩阵的范围
        position = dict()
        for r in range(row):
            for c in range(col):
                color = targetGrid[r][c]
                if color not in position:
                    position[color] = [r, r, c, c]  # i最小值，i最大值，j最小值，j最大值
                else:
                    position[color] = [
                        min(r, position[color][0]),
                        max(r, position[color][1]),
                        min(c, position[color][2]),
                        max(c, position[color][3]),
                    ]

        # 依赖关系
        # 遍历每个数字的矩阵范围内包含的其他数字
        adjMap = defaultdict(set)
        indegree = {color: 0 for color in position.keys()}  ## 注意不能用defaultdict 因为要包含所有点
        visited = set()
        for color1, pos in position.items():
            for r in range(pos[0], pos[1] + 1):
                for c in range(pos[2], pos[3] + 1):
                    color2 = targetGrid[r][c]
                    if color2 != color1:
                        adjMap[color1].add(color2)
                        if (color1, color2) not in visited:
                            visited.add((color1, color2))
                            indegree[color2] += 1

        # 拓扑排序环检测即可
        queue = deque([c for c, d in indegree.items() if d == 0])
        count = 0
        while queue:
            color = queue.popleft()
            count += 1
            for next in adjMap[color]:
                indegree[next] -= 1
                if indegree[next] == 0:
                    queue.append(next)

        return len(position) == count


print(Solution().isPrintable(targetGrid=[[1, 1, 1, 1], [1, 1, 3, 3], [1, 1, 3, 4], [5, 5, 1, 4]]))
# 输出：true


# 例如
# [1,1,1,1]
# [1,1,3,3]
# [1,1,3,4]
# [5,5,1,4]

# 在这个例子中，
# 1的矩形包含3，4，5 =>  1要在他们前面
# 3的矩形包含4    =>    3要在4前面
