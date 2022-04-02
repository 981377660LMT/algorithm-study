# n组命题，表示区间[a,b]里单词的个数是奇数还是偶数
# 问命题是否冲突
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


class Solution:
    def solve(self, lists):
        n = len(lists)
        # 开大一点
        uf = UnionFindArray(4 * n + 10)

        for a, b, type in lists:
            b = b + 1
            # 第一个冲突
            if type == 0:
                # 相同的不能在一组
                if uf.isConnected(a, b + 2 * n):
                    return False
                uf.union(a, b)
                uf.union(a + 2 * n, b + 2 * n)
            else:
                # 不同的在一组了
                if uf.isConnected(a, b):
                    return False
                uf.union(a, b + 2 * n)
                uf.union(a + 2 * n, b)
        return True


# print(Solution().solve(lists=[[1, 5, 1], [6, 10, 0], [1, 10, 0]]))
# If there are an odd number of words from pages [1, 5] and an even number of words from pages [6, 10],
# that must mean there's an odd number from pages [1, 10]. So this is a contradiction.
print(Solution().solve(lists=[[0, 2, 0], [2, 2, 0]]))
