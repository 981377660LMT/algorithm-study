from functools import reduce
from typing import List
from collections import Counter

# # 总结：
# # 简单说就是(root - currentNode) * currentNode.left * currentNode.right
# # 后序dfs统计left/right的字树数量


class Solution:
    def countHighestScoreNodes(self, parents: List[int]) -> int:
        def dfs(cur: int, parent: int) -> int:
            """dfs后序返回子树结点数"""
            subTree = []
            for next in adjList[cur]:
                if next == parent:
                    continue
                subTree.append(dfs(next, cur))

            sum_ = sum(subTree)
            parentCount = max(1, n - sum_ - 1)
            score = parentCount * reduce(lambda x, y: x * y, subTree, 1)
            counter[score] += 1
            return sum_ + 1

        n = len(parents)
        adjList = [[] for _ in range(n)]
        for i, parent in enumerate(parents):
            if parent != -1:
                adjList[parent].append(i)
                adjList[i].append(parent)

        counter = Counter()
        dfs(0, -1)
        return counter[max(counter)]


# print(Solution().countHighestScoreNodes(parents=[-1, 2, 0, 2, 0]))
# # 输出：3
# # 解释：
# # - 节点 0 的分数为：3 * 1 = 3
# # - 节点 1 的分数为：4 = 4
# # - 节点 2 的分数为：1 * 1 * 2 = 2
# # - 节点 3 的分数为：4 = 4
# # - 节点 4 的分数为：4 = 4
# # 最高得分为 4 ，有三个节点得分为 4 （分别是节点 1，3 和 4 ）。
