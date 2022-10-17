# 现给出一个N行M列矩阵(N,M<=2e5)，并给出Q次询问，如下:
# 1.对列给[left,right]所有列元素加上x  形如 1 left right x
# 2.对行将所有第i行元素赋值为x  形如 2 i x
# 3.输出第i行j列元素  形如 3 i j

# !离线查询,预处理query
# https://www.zhihu.com/search?type=content&q=AtCoder%20Beginner%20Contest%20253%20
# 考虑没有操作2的情况，那么很容易地就可以用树状数组实现对列的区间加及单点查询
# !维护列的前缀和 记录每行上一次修改是在什么时候
# !遍历到这个时候就减去之前这个点的值再加上现在赋值的值

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


from collections import defaultdict


class BIT2:
    """范围修改"""

    def __init__(self, n: int):
        self.size = n
        self._tree1 = defaultdict(int)
        self._tree2 = defaultdict(int)

    def add(self, left: int, right: int, delta: int) -> None:
        """闭区间[left, right]加delta"""
        self._add(left, delta)
        self._add(right + 1, -delta)

    def query(self, left: int, right: int) -> int:
        """闭区间[left, right]的和"""
        return self._query(right) - self._query(left - 1)

    def _add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")

        rawIndex = index
        while index <= self.size:
            self._tree1[index] += delta
            self._tree2[index] += (rawIndex - 1) * delta
            index += index & -index

    def _query(self, index: int) -> int:
        if index > self.size:
            index = self.size

        rawIndex = index
        res = 0
        while index > 0:
            res += rawIndex * self._tree1[index] - self._tree2[index]
            index -= index & -index
        return res


if __name__ == "__main__":
    ROW, COL, q = map(int, input().split())
    colBit = BIT2(int(2e5 + 10))  # 列的前缀和
    lastUpdate = [0] * (ROW + 10)  # 记录每行上一次赋值(操作2)是在什么时候
    toUpdateCol = [[] for _ in range(q + 10)]  # 需要在这个时刻更新的列
    res = [0] * (q + 10)
    queries = []
    queries.append(())  # 占位

    for i in range(1, q + 1):
        op, *rest = map(int, input().split())
        if op == 1:
            left, right, delta = rest
            queries.append((op, left, right, delta))
        elif op == 2:
            row, target = rest
            lastUpdate[row] = i
            queries.append((op, row, target))
        else:
            row, col = rest
            toUpdateCol[lastUpdate[row]].append((i, col))  # 记录需要在这个时刻更新的列
            queries.append((op, row, col))

    for i in range(1, q + 1):
        op, *rest = queries[i]
        if op == 1:
            left, right, delta = rest
            colBit.add(left, right, delta)
        elif op == 2:
            row, target = rest
            for qi, col in toUpdateCol[i]:
                res[qi] -= colBit.query(col, col)
                res[qi] += target
        else:
            row, col = rest
            res[i] += colBit.query(col, col)
            print(res[i])
