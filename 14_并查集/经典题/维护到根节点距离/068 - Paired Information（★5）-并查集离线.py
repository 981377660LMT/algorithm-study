# q个条件
# !Ti=0 时 给出数组相邻两项的和 nums[i]+nums[i+1]=vi
# !TI=1 时 假定nums[xi]=Vi 判断nums[yi]是否有确定的值
# 题目不会给出矛盾的数据
# 如果不确定 输出 'Ambiguous'
# N,Q<=1e5
# Xi,Yi<=2e9

# 预处理查询 把连通性和各个约束大小处理完

import sys


sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加
    
    初始化的连通分量个数 为 n
    """

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))

    def find(self, key: int) -> int:
        while self.parent[key] != key:
            self.parent[key] = self.parent[self.parent[key]]
            key = self.parent[key]
        return key

    def union(self, x: int, y: int) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        self.parent[rootX] = rootY
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


n = int(input())
q = int(input())
queries = []
for _ in range(q):
    t, i, j, v = map(int, input().split())  # i,j>=1
    i, j = i - 1, j - 1
    queries.append((t, i, j, v))

# 预处理查询 クエリ先読み
pairSum = [0] * (n - 1)  # s[i] = A[i] + A[i+1]
for t, i, j, v in queries:
    if t == 0:
        pairSum[i] = v

# 处理偏移距离(开头元素基准为0)
offsets = [0] * n
for i in range(n - 1):
    offsets[i + 1] = pairSum[i] - offsets[i]

uf = UnionFindArray(n + 5)
for t, i, j, v in queries:
    if t == 0:
        uf.union(i, j)
    else:
        if not uf.isConnected(i, j):
            print('Ambiguous')
        else:
            res = 0
            # !A[i]がx増えると，A[i+1]はx減り，A[i+2]はx増え...というように変化する
            if (i + j) % 2 == 0:
                res = offsets[j] + (v - offsets[i])
            else:
                res = offsets[j] - (v - offsets[i])
            print(res)
