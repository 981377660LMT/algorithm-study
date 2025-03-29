# abc399-E - Replace(置换环)
# https://atcoder.jp/contests/abc399/tasks/abc399_e
#
# 给定正整数 N 以及长度为 N 的由英文字母小写字母组成的字符串 S 和 T。
# 允许执行如下操作任意次（可为 0 次）：
# !选择任意两个英小文字 x 和 y，将 S 中所有的 x 同时替换为 y。
# !请判断是否可以通过若干次该操作使 S 与 T 完全一致；如果可以，请求出所需操作次数的最小值。
#
# 因为每个点的出度一定为1，所以图最多只有1个环，则图由以下6种情况:
# A类：孤立顶点
# B类：长度为1的环
# C类：长度≥2的环
# D类：有根树
# E类：长度1的环+附着的有根树
# F类：长度≥2的环+附着的有根树


from typing import List


class UnionFindArraySimple:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(self, key1: int, key2: int) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        if self._data[key] < 0:
            return key
        self._data[key] = self.find(self._data[key])
        return self._data[key]

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]

    def getGroups(self) -> List[List[int]]:
        res = [[] for _ in range(self.n)]
        for i in range(self.n):
            res[self.find(i)].append(i)
        return [x for x in res if x]


C = 26
if __name__ == "__main__":
    N = int(input())
    S = input()
    T = input()

    if S == T:
        print(0)
        exit()

    nums1 = [ord(c) - ord("a") for c in S]
    nums2 = [ord(c) - ord("a") for c in T]
    to = [-1] * C
    for c1, c2 in zip(nums1, nums2):
        if to[c1] != -1 and to[c1] != c2:
            print(-1)
            exit()
        to[c1] = c2

    tmp = sorted(to)
    isPerm = all(tmp[i] == i for i in range(C))
    if isPerm:
        print(-1)
        exit()

    res = 0
    indeg = [0] * C
    uf = UnionFindArraySimple(C)
    for i, v in enumerate(to):
        if v != -1:
            if i != v:
                res += 1
            indeg[v] += 1
            uf.union(i, v)

    for g in uf.getGroups():
        if len(g) == 1:
            continue
        isCycle = all(indeg[i] == 1 for i in g)
        res += isCycle

    print(res)
