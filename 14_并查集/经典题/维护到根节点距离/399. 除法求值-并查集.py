# https://leetcode.cn/problems/evaluate-division/submissions/

from typing import List
from UnionFindWithDist import UnionFindMapWithDist2


class Solution:
    def calcEquation(
        self, equations: List[List[str]], values: List[float], queries: List[List[str]]
    ) -> List[float]:
        """如果存在某个无法确定的答案，则用 -1.0 替代这个答案。
        如果问题中出现了给定的已知条件中没有出现的字符串，也需要用 -1.0 替代这个答案。
        乘积关系取对数就是加法 等价于维护到根节点的距离
        """

        def id(o: object) -> int:
            if o not in _pool:
                _pool[o] = len(_pool)
            return _pool[o]

        _pool = dict()

        uf = UnionFindMapWithDist2()
        for (key1, key2), value in zip(equations, values):
            uf.union(id(key1), id(key2), value)  # !key1/key2(距离)为 value
        res = []
        for u, v in queries:
            id1, id2 = id(u), id(v)
            if id1 not in uf or id2 not in uf or not uf.isConnected(id1, id2):
                res.append(-1.0)
            else:
                res.append(uf.dist(id1, id2))

        return res


print(
    Solution().calcEquation(
        [["a", "b"], ["b", "c"]],
        [2.0, 3.0],
        [["a", "c"], ["b", "a"], ["a", "e"], ["a", "a"], ["x", "x"]],
    )
)
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
