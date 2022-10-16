# https://leetcode.cn/problems/evaluate-division/submissions/

from typing import List
from UnionFindMapWithDist import UnionFindMapWithDist1


class Solution:
    def calcEquation(
        self, equations: List[List[str]], values: List[float], queries: List[List[str]]
    ) -> List[float]:
        """如果存在某个无法确定的答案，则用 -1.0 替代这个答案。
        如果问题中出现了给定的已知条件中没有出现的字符串，也需要用 -1.0 替代这个答案。
        乘积关系取对数就是加法 等价于维护到根节点的距离
        """
        uf = UnionFindMapWithDist1[str]()
        for (key1, key2), value in zip(equations, values):
            uf.add(key1).add(key2).union(key2, key1, value)  # !value * key2 = key1

        res = []
        for u, v in queries:
            if u not in uf or v not in uf or not uf.isConnected(u, v):
                res.append(-1.0)
            else:
                res.append(uf.distToRoot[v] / uf.distToRoot[u])

        return res


# print(
#     Solution().calcEquation(
#         [["a", "b"], ["b", "c"]],
#         [2.0, 3.0],
#         [["a", "c"], ["b", "a"], ["a", "e"], ["a", "a"], ["x", "x"]],
#     )
# )
# print(
#     Solution().calcEquation(
#         equations=[["a", "b"], ["b", "c"], ["bc", "cd"]],
#         values=[1.5, 2.5, 5.0],
#         queries=[["a", "c"], ["c", "b"], ["bc", "cd"], ["cd", "bc"]],
#     )
# )
print(
    Solution().calcEquation(
        equations=[["a", "e"], ["b", "e"]],
        values=[4.0, 3.0],
        queries=[["a", "b"], ["e", "e"], ["x", "x"]],
    )
)
# [1.33333,1.0,-1.0]
