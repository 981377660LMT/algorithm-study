from typing import List
from collections import defaultdict, deque

# 每一次操作时，打印机会用同一种颜色打印一个矩形的形状，每次打印会覆盖矩形对应格子里原本的颜色。
# 一旦矩形根据上面的规则使用了一种颜色，那么` 相同的颜色不能再被使用 。`

# 给你一个初始没有颜色的 m x n 的矩形 targetGrid
# 如果你能按照上述规则打印出矩形 targetGrid ，请你返回 true ，否则返回 false 。

# 单向图不允许形成环


class Solution:
    def isPrintable(self, targetGrid: List[List[int]]) -> bool:
        # 遍历所有点，计算出每个值矩阵的范围
        dic = {}
        for i in range(len(targetGrid)):
            for j in range(len(targetGrid[0])):
                cur = targetGrid[i][j]
                if not cur in dic:
                    dic[cur] = [i, i, j, j]  # i最小值，i最大值，j最小值，j最大值
                else:
                    dic[cur] = [
                        min(i, dic[cur][0]),
                        max(i, dic[cur][1]),
                        min(j, dic[cur][2]),
                        max(j, dic[cur][3]),
                    ]

        # 依赖关系
        # 遍历每个数字的矩阵范围内包含的其他数字
        dependency = defaultdict(set)
        indegree = {color: 0 for color in dic.keys()}
        visited = set()
        for k, v in dic.items():
            for i in range(v[0], v[1] + 1):
                for j in range(v[2], v[3] + 1):
                    if targetGrid[i][j] != k:
                        dependency[k].add(targetGrid[i][j])
                        if (k, targetGrid[i][j]) not in visited:
                            visited.add((k, targetGrid[i][j]))
                            indegree[targetGrid[i][j]] += 1

        # 拓扑排序环检测即可
        queue = deque([c for c, v in indegree.items() if v == 0])
        target = len(dic)
        count = 0
        while queue:
            cur = queue.popleft()
            count += 1
            for next in dependency[cur]:
                indegree[next] -= 1
                if indegree[next] == 0:
                    queue.append(next)

        return target == count


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
