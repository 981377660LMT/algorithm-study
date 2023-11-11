# 维护一个初始是全0的 n x m (n, m ≤2e5）的矩阵，支持:
# 1.对[Li,ri]这些列的每一个元素加x
# 2.将第i行全部赋值为x
# 3.查询矩阵中(xi , yi）的值


# !在线的做法就是写一个主席树 + 标记持久化；
# !离线的做法就是把后面的贡献写做前缀和差分，然后两个时刻维护一下。
####################################################################################
# 离线做法
# 对于每一行，每一列单独维护
# !可以計算出第j列累加的值s，第i行最後一次變的值x，還有變成x的時候第j列累加的值s1，結果就是x+s-s1，使用離線做法
# 涉及到染色 可以倒着处理查询

from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

# TODO


def main() -> None:
    ROW, COL, q = map(int, input().split())
    Q = [tuple(map(int, input().split())) for _ in range(q)]

    rowUpdate = defaultdict(lambda: [0])  # 每个行的更新时刻
    bit = BIT2(int(3e5))  # 维护列
    q3 = defaultdict(set)
    res = []

    for qt, *rest in Q[::-1]:
        if qt == 1:
            left, right, delta = rest
            bit.add(left, right, delta)
        elif qt == 2:
            rowUpdate[rest[0]].append(rest[1])
            bit.add(rest[1], -1)
        else:
            res.append(bit.query(rest[0], rest[1]))
            q3[rest[0]].add(rest[1])
            for i in q3[rest[0]]:
                bit.add(i, -1)


if __name__ == "__main__":

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

        if os.environ.get("USERNAME", " ") == "caomeinaixi":
            while True:
                main()
        else:
            main()
