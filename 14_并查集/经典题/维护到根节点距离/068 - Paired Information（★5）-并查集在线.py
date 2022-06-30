# q个条件
# !Ti=0 时 给出数组相邻两项的和 nums[i]+nums[i+1]=vi
# !TI=1 时 假定nums[xi]=Vi 判断nums[yi]是否有确定的值
# 题目不会给出矛盾的数据
# 如果不确定 输出 'Ambiguous'
# N,Q<=1e5
# Xi,Yi<=2e9

# !技巧:把相邻两项和转化为相邻两项差 从而可以用BIT/Seg快速计算任意两项之间的差
# !nums[i]+nums[i+1]=vi 变为 nums[i+1]-nums[i]=-(-1)^i*vi
# !nums[xi]=Vi 变为 nums[xi]=(-1)^i*Vi


import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


class DistUnionFindArray:
    """在并查集里面维护到根节点的距离"""

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.distToRoot = [0] * n

    def find(self, key: int) -> int:
        if key != self.parent[key]:
            root = self.find(self.parent[key])
            # !x到根节点的距离更新了
            self.distToRoot[key] += self.distToRoot[self.parent[key]]
            self.parent[key] = root
        return self.parent[key]

    def union(self, x: int, y: int, diff: int) -> bool:
        """认key2作为key1的父节点"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        self.parent[rootX] = rootY
        self.distToRoot[rootY] += self.distToRoot[rootX] + diff
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


# n = int(input())
# q = int(input())
# uf = DistUnionFindArray(n + 10)

# for _ in range(q):
#     t, i, j, v = map(int, input().split())
#     i, j = i - 1, j - 1
#     if t == 0:
#         uf.union(i, j, pow(-1, i) * v)
#     else:
#         if i == j:
#             print(v)
#         if not uf.isConnected(i, j):
#             print('Ambiguous')
#         else:
#             res = 0
#             if i < j:
#                 diff = uf.distToRoot[i] - uf.distToRoot[j]
#                 res = diff * pow(-1, j + 1) + pow(-1, i + j) * v
#             else:
#                 diff = uf.distToRoot[j] - uf.distToRoot[i]
#                 res = diff * pow(-1, j + 1) + pow(-1, i + j) * v
#             print(res)

# 哪里不太对
