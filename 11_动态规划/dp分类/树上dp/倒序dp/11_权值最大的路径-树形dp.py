# 多源最短路?
# floyd

from collections import defaultdict
from typing import List


INF = 0x3F3F3F3F
# 给出一个有向无环图，规定路径是单向且`小序号指向大序号`，每个节点都有权值。
# 在图上求一条路径使得经过的节点权值和最大，输出路径

# n<=10^4
# 树形dp 倒序从后往前


class Solution:
    def digSum(self, potatoNum: List[int], connectRoad: List[List[int]]) -> str:
        n = len(potatoNum)
        adjList = defaultdict(list)
        score, path = [0] * (n + 1), [''] * (n + 1)
        for cur, next in connectRoad:
            adjList[cur].append(next)

        # 树形倒序dp
        for root in range(n, 0, -1):
            score[root] = potatoNum[root - 1]
            path[root] = str(root)
            maxScore, nextCand = -1, -1
            for next in adjList[root]:
                if score[next] > maxScore:
                    maxScore = score[next]
                    nextCand = next
            if maxScore != -1:
                score[root] += maxScore
                path[root] = path[root] + '-' + path[nextCand]

        maxScore, maxPath = -1, ''
        for i in range(1, n + 1):
            if score[i] > maxScore:
                maxScore = score[i]
                maxPath = path[i]
        return maxPath


print(
    Solution().digSum(
        [5, 10, 20, 5, 4, 5], [[1, 2], [1, 4], [2, 4], [3, 4], [4, 5], [4, 6], [5, 6]]
    )
)

