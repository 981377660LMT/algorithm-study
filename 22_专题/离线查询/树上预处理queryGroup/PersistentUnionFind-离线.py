# https://www.cnblogs.com/antiquality/p/9427175.html
# 有m次操作,操作分为3种
# !1 a b 合并a,b所在集合 (1<=a,b<=n)
# !2 k 回到`第k次操作`之后的状态(k>=0)
# !3 a b 查询a,b是否在同一个集合 如果是输出1 否则输出0

# 把版本看成顶点 那么所有版本构成了一棵树
# !预处理查询组后在树上dfs输出所有查询结果(因为要回溯,所以使用可撤销并查集)

from typing import List


def solve(n: int, operations: List[List[int]]) -> List[int]:
    def dfs(curVersion: int) -> None:
        """版本树上dfs处理所有查询"""
        for num1, num2, qi in queryGroup[curVersion]:
            res[qi] = 1 if uf.isConnected(num1, num2) else 0
        for next, num1, num2 in adjList[curVersion]:
            uf.union(num1, num2)
            dfs(next)
            uf.revocate()

    m = len(operations)
    uf = RevocableUnionFindArray(n + 10)
    adjList = [[] for _ in range(m + 10)]  # !版本间转移形成的树
    queryGroup = [[] for _ in range(m + 10)]  # !每个版本处的查询
    git = [0] * (m + 10)  # !每次操作后的版本号,初始版本为0
    curVersion = 0
    for i in range(m):
        kind, *args = operations[i]
        if kind == 1:  # action
            num1, nums2 = args
            adjList[curVersion].append((curVersion + 1, num1, nums2))
            curVersion += 1
        elif kind == 2:  # action
            k = args[0]
            preVersion = git[k]
            adjList[preVersion].append((curVersion + 1, 0, 0))
            curVersion += 1
        elif kind == 3:  # query
            num1, nums2 = args
            queryGroup[curVersion].append((num1, nums2, i))
        git[i + 1] = curVersion

    res = [-1] * (m + 10)
    dfs(0)
    return [num for num in res if num != -1]


class RevocableUnionFindArray:
    """
    带撤销操作的并查集

    不能使用路径压缩优化（因为路径压缩会改变结构）；
    为了不超时必须使用按秩合并优化,复杂度nlogn
    """

    __slots__ = ("n", "part", "parent", "rank", "optStack")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n
        self.optStack = []

    def find(self, x: int) -> int:
        """不能使用路径压缩优化"""
        while self.parent[x] != x:
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        """x所在组合并到y所在组"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            self.optStack.append((-1, -1, -1))
            return False

        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX

        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
        self.optStack.append((rootX, rootY, self.rank[rootX]))
        return True

    def revocate(self) -> None:
        """
        用一个栈记录前面的合并操作，
        撤销时要依次取出栈顶元素做合并操作的逆操作
        """
        if not self.optStack:
            raise IndexError("no union option to revocate")

        rootX, rootY, rankX = self.optStack.pop()
        if rootX == -1:
            return

        self.parent[rootX] = rootX
        self.rank[rootY] -= rankX
        self.part += 1

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, m = map(int, input().split())
    operations = []
    for _ in range(m):
        operations.append(list(map(int, input().split())))
    res = solve(n, operations)
    for num in res:
        print(num)
