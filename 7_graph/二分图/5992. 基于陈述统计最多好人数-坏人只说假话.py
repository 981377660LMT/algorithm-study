from typing import List, Tuple
from collections import defaultdict

# 1 <= 人数 <= 2*10^5，1 <= 评论条数 (边数) <= 5*10^5
# 题目变成了好人只说真话，坏人只说假话
# 如果i说j是坏人,ij身份不同；如果i说j是好人,ij身份相同
# 二分图
# 如果i说j坏人，ij直接连边(身份不同)；否则设置一个中间点fake(编号为n+1)将ij连起来(身份相同)
# 对每一个强连通分量做二分图判定即可；即对每个未染色的点出发dfs做二分图判定


class Solution:
    def maximumGood(self, statements: List[List[int]]) -> int:
        def dfs(cur: int, color: int, colorCounter: List[int]) -> bool:
            """"二分图检测，并统计两种颜色的数量"""
            colors[cur] = color
            # 计数时忽略虚拟结点
            if cur < n:
                colorCounter[color] += 1
            for next in adjMap[cur]:
                if colors[next] == -1:
                    if not dfs(next, color ^ 1, colorCounter):
                        return False
                elif colors[cur] == colors[next]:
                    return False
            return True

        # 建图
        n = len(statements)
        adjMap = defaultdict(set)
        dummyId = n
        for curId, row in enumerate(statements):
            for nextId, state in enumerate(row):
                if state == 0:
                    adjMap[curId].add(nextId)
                    adjMap[nextId].add(curId)
                elif state == 1:
                    adjMap[curId].add(dummyId)
                    adjMap[dummyId].add(curId)
                    adjMap[nextId].add(dummyId)
                    adjMap[dummyId].add(nextId)
                    dummyId += 1

        # 二分图检测，统计颜色
        res = 0
        colors = defaultdict(lambda: -1)
        for i in range(n):
            if colors[i] != -1:
                continue
            colorCounter = [0, 0]
            isBiPartite = dfs(i, 0, colorCounter)
            if not isBiPartite:
                return -1
            res += max(colorCounter)

        return res


print(Solution().maximumGood(statements=[[2, 1, 2], [1, 2, 2], [2, 0, 2]]))

