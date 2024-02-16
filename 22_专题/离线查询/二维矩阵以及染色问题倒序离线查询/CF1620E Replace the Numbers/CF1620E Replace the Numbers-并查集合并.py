# CF1620E Replace the Numbers-在线查询
# https://www.luogu.com.cn/problem/CF1620E
# 给出 q 个操作，操作分为两种：

# 1 x 在序列末尾插入数字 x。
# 2 x y 把序列中的所有 x 替换为 y。

# 求这个序列操作后的结果。
# 并查集，类似未来日记中的技巧

import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")


class UfMap:
    __slots__ = "_parent"

    def __init__(self):
        self._parent = dict()

    def unionTo(self, child: int, parent: int) -> bool:
        root1, root2 = self.find(child), self.find(parent)
        if root1 == root2:
            return False
        self._parent[root1] = root2
        return True

    def find(self, key: int) -> int:
        if key not in self._parent:
            self._parent[key] = key
            return key
        while self._parent[key] != key:
            self._parent[key] = self._parent[self._parent[key]]
            key = self._parent[key]
        return key


class UnionFindMapSimple:
    __slots__ = "_data"

    def __init__(self):
        self._data = dict()

    def add(self, key: int) -> None:
        self._data[key] = -1

    def union(self, parent: int, child: int) -> bool:
        root1, root2 = self.find(parent), self.find(child)
        if root1 == root2:
            return False
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        return True

    def find(self, key: int) -> int:
        if key not in self._data:
            self.add(key)
            return key
        if self._data[key] < 0:
            return key
        self._data[key] = self.find(self._data[key])
        return self._data[key]


if __name__ == "__main__":
    q = int(input())
    firstPos = dict()  # 每个值出现的第一个位置
    values = []  # 每个位置对应的值
    uf = UnionFindMapSimple()  # 维护哪些位置的值相同
    count = 0

    def add(index: int, value: int) -> None:
        values.append(value)
        if value not in firstPos:
            firstPos[value] = index
        else:
            uf.union(index, firstPos[value])

    def merge(from_: int, to: int) -> None:
        if from_ == to:
            return
        if from_ not in firstPos:
            return
        if to not in firstPos:
            firstPos[to] = firstPos[from_]
            values[firstPos[to]] = to
            firstPos.pop(from_)
            return
        uf.union(firstPos[from_], firstPos[to])
        firstPos.pop(from_)

    for _ in range(q):
        t, *args = map(int, input().split())
        if t == 1:
            x = args[0]
            add(count, x)
            count += 1
        else:
            x, y = args
            merge(y, x)

    res = [values[uf.find(i)] for i in range(count)]
    print(*res)
