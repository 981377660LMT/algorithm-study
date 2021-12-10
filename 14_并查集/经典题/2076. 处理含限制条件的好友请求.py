from typing import List
from collections import defaultdict

# 2 <= n <= 1000

# restrictions[i] = [xi, yi] 意味着用户 xi 和用户 yi 不能 成为 朋友

# 如果第 j 个好友请求 成功 ，那么 result[j] 就是 true
class UnionFind:
    def __init__(self, n):
        self.parent = list(range(n))

    def union(self, x, y):
        rx, ry = self.find(x), self.find(y)
        if rx == ry:
            return False
        low, high = sorted([rx, ry])
        self.parent[high] = low
        return True

        # 关键在于合并forbidden

    def find(self, i):
        if i != self.parent[i]:
            self.parent[i] = self.find(self.parent[i])
        return self.parent[i]


class Solution:
    def friendRequests(
        self, n: int, restrictions: List[List[int]], requests: List[List[int]]
    ) -> List[bool]:
        res = [True] * len(requests)
        uf = UnionFind(n)

        for index, [u, v] in enumerate(requests):
            isSuccess = True
            ru, rv = uf.find(u), uf.find(v)

            if ru != rv:
                for x, y in restrictions:
                    # 这一对不能合并
                    rx, ry = uf.find(x), uf.find(y)
                    if (rx, ry) == (ru, rv) or (rx, ry) == (rv, ru):
                        isSuccess = False
                        break

            res[index] = isSuccess
            if isSuccess:
                uf.union(u, v)

        return res


print(Solution().friendRequests(n=3, restrictions=[[0, 1]], requests=[[0, 2], [2, 1]]))
# 输出：[true,false]
# 解释：
# 请求 0 ：用户 0 和 用户 2 可以成为朋友，所以他们成为直接朋友。
# 请求 1 ：用户 2 和 用户 1 不能成为朋友，因为这会使 用户 0 和 用户 1 成为间接朋友 (1--2--0) 。

