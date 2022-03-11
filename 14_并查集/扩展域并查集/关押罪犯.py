# 他准备将罪犯们在两座监狱内重新分配，\
# 以求产生的冲突事件影响力都较小，
# 从而保住自己的乌纱帽。
# 假设只要处于同一监狱内的某两个罪犯间有仇恨，
# 那么他们一定会在每年的某个时候发生摩擦。
# 那么，应如何分配罪犯，才能使Z市长看到的那个冲突事件的影响力最小？这个最小值是多少？
# 公务繁忙的Z 市长只会去看列表中的第一个事件的影响力，如果影响很坏，他就会考虑撤换警察局长。
from typing import List


class UnionFindArray:
    def __init__(self, n: int):
        self.n = n
        self.count = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.count -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


# 互相不喜欢的人可以用特殊的并查集，每个集分为 2 部分，不喜欢的人一定在同一个集的另一部分里面
class Solution:
    def 关押罪犯(self, n: int, dislikes: List[List[int]]) -> int:
        uf = UnionFindArray(n * 2 + 10)
        dislikes.sort(key=lambda x: x[2], reverse=True)
        for cur, next, score in dislikes:
            # 第一个冲突
            if uf.isConnected(cur, next):
                return score
            uf.union(cur, next + n)
            uf.union(cur + n, next)
        return 0

