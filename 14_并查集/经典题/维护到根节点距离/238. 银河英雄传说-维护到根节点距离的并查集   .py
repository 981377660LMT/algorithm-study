# 有 T 条指令，每条指令格式为以下两种之一：
# M i j，表示让第 i 号战舰所在列的全部战舰保持原有顺序，接在第 j 号战舰所在列的尾部。
# C i j，表示询问第 i 号战舰与第 j 号战舰当前是否处于同一列中，如果在同一列中，它们之间间隔了多少艘战舰。


class UnionFindArrayWithDist:
    """固定元素 维护距离的并查集"""

    def __init__(self, n: int):
        self.n = n
        self.count = n
        self.parent = list(range(n))
        self.rank = [1] * n
        self.distToRoot = [0] * n

    def find(self, x: int) -> int:
        """注意有路径压缩 除了第一次调用 之后的distToRoot不会继续变化"""
        if x != self.parent[x]:
            root = self.find(self.parent[x])
            # !x到根节点的距离更新了
            self.distToRoot[x] += self.distToRoot[self.parent[x]]
            self.parent[x] = root
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        """x接到y上"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        # if self.rank[rootX] > self.rank[rootY]:
        #     rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        # !注意这里 距离增加为帮派大小
        self.distToRoot[rootX] += self.rank[rootY]
        self.rank[rootY] += self.rank[rootX]
        self.count -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


# 并查集+懒更新
# N≤30000,T≤500000


def main():
    uf = UnionFindArrayWithDist(30000 + 10)

    T = int(input())
    for _ in range(T):
        line = list(input().split())
        op = line[0]
        x, y = int(line[1]), int(line[2])

        if op == "M":
            uf.union(x, y)
        else:
            if not uf.isConnected(x, y):
                print(-1)
            else:
                cur = abs(uf.distToRoot[x] - uf.distToRoot[y])
                print(max(0, cur - 1))


if __name__ == "__main__":
    main()
