from collections import defaultdict


# 判断一个有向图的是否存在欧拉路径
# 1. 所有点连通（并查集）
# 2. 只有两个入度不等于出度的点 或者所有点都是出度等于入度
#    1. 起点 出度 比 入度 多一
#    2. 终点 入度 比 出度 多一

# 有 N 个盘子，每个盘子上写着一个仅由小写字母组成的英文单词。
# 你需要给这些盘子安排一个合适的顺序，使得相邻两个盘子中，前一个盘子上单词的末字母等于后一个盘子上单词的首字母。
# 如果存在合法解，则输出”Ordering is possible.”，否则输出”The door cannot be opened.”。

# 判断是否存在欧拉路径 s0=>a=>s2=>b=>s3=>...=>s26=>z=>sn


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
    def check() -> bool:
        # 判断有向图连通性
        pre = -1
        for cur in range(26):
            if not visited[cur]:
                continue
            if pre == -1:
                pre = cur
                continue
            if not uf.isConnected(pre, cur):
                return False
            pre = cur

        # 判断是否存在有向图欧拉路径
        startCount, endCount = 0, 0
        for i in range(26):
            if not visited[i]:
                continue
            diff = outd[i] - ind[i]
            if diff == 0:
                continue
            elif diff == 1:
                startCount += 1
            elif diff == -1:
                endCount += 1
            else:
                return False
        return (startCount, endCount) in [(0, 0), (1, 1)]

    n = int(input())
    ind, outd = [0] * 26, [0] * 26
    visited = [0] * 26
    uf = UnionFindArray(26 + 10)
    for _ in range(n):
        word = input()
        first, last = ord(word[0]) - 97, ord(word[-1]) - 97
        visited[first] = visited[last] = True
        outd[first] += 1
        ind[last] += 1
        uf.union(first, last)

    if check():
        print("Ordering is possible.")
    else:
        print("The door cannot be opened.")


t = int(input())
for _ in range(t):
    main()

