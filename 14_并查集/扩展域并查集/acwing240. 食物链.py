# 动物王国中有三类动物 A,B,C，这三类动物的食物链构成了有趣的环形。
# A 吃 B，B 吃 C，C 吃 A。
# 现有 N 个动物，以 1∼N 编号。
# 每个动物都是 A,B,C 中的一种，但是我们并不知道它到底是哪一种。
# 有人用两种说法对这 N 个动物所构成的食物链关系进行描述：
# 第一种说法是 1 X Y，表示 X 和 Y 是同类。
# 第二种说法是 2 X Y，表示 X 吃 Y。
# 此人对 N 个动物，用上述两种说法，一句接一句地说出 K 句话，这 K 句话有的是真的，有的是假的。
# 当一句话满足下列三条之一时，这句话就是假话，否则就是真话。

# 当前的话与前面的某些真的话冲突，就是假话；
# 当前的话中 X 或 Y 比 N 大，就是假话；
# 当前的话表示 X 吃 X，就是假话。
# 你的任务是根据给定的 N 和 K 句话，输出假话的总数。


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


def main():
    n, k = map(int, input().split())
    uf = UnionFindArray(3 * n + 10)
    res = 0

    for _ in range(k):
        f, x, y = map(int, input().split())
        if x > n or y > n:
            res += 1
            continue
        # ------- x和y是同类
        if f == 1:
            if uf.isConnected(x, y + n) or uf.isConnected(x, y + 2 * n):
                res += 1
            else:
                uf.union(x, y)  # 同类域
                uf.union(x + n, y + n)  # 吃域
                uf.union(x + 2 * n, y + 2 * n)  # 被吃域

        # ------- x吃y
        elif f == 2:
            if x == y:
                res += 1
                continue
            if uf.isConnected(x, y) or uf.isConnected(x, y + n):
                res += 1
            else:
                uf.union(x + n, y)  # x的吃域与y同类
                uf.union(x + 2 * n, y + n)  # x吃域的吃域与y的吃域同类
                uf.union(x, y + 2 * n)  # y吃域的吃域与x同类
    print(res)


if __name__ == '__main__':
    main()

