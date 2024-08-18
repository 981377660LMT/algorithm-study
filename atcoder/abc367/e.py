from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 各要素が
# 1 以上
# N 以下である長さ
# N の数列
# X と、長さ
# N の数列
# A が与えられます。
# A に以下の操作を
# K 回行った結果を出力してください。

# B
# i
# ​
#  =A
# X
# i
# ​

# ​
#   なる
# B を新たな
# A とする


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


if __name__ == "__main__":
    N, K = map(int, input().split())
    X = list(map(int, input().split()))
    A = list(map(int, input().split()))
    for i in range(N):
        X[i] -= 1
        A[i] -= 1
    uf = UnionFindArraySimple(N)
    for i in range(N):
        uf.union(i, X[i])

    groups = defaultdict(list)
    for i in range(N):
        root = uf.find(i)
        groups[root].append(i)

    print(groups)
    res = [0] * N
    for group in groups.values():
        size = len(group)
        tmpGroup = group[:]
        for i in range(size):
            res[group[i]] = tmpGroup[(i + K) % size]
    for i in range(N):
        print(A[res[i]] + 1, end=" ")
