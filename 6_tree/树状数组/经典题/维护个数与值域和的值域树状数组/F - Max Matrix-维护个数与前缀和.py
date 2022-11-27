# https://atcoder.jp/contests/jsc2021/tasks/jsc2021_f

# 输入 n m q (≤2e5)，初始你有长为 n 的数组 a，长为 m 的数组 b，元素值都为 0，下标从 1 开始。
# 然后输入 q 个询问，每个询问形如 t x y (1≤y≤1e8)。
# t=1，表示把 a[x]=y；t=2，表示把 b[x]=y。
# 每次修改后，输出 ∑∑max(a[i],b[j])，这里 i 取遍 [1,n]，j 取遍 [1,m]。

# https://atcoder.jp/contests/jsc2021/submissions/36189511
# https://atcoder.jp/contests/jsc2021/submissions/36189772

# 提示 1：元素的顺序并没有那么重要。

# !提示 2：每次修改时，只需要知道另一个数组中，有多少个数比自己小（这些数的个数乘上自己），以及≥自己的数的和是多少。

# 最好写的应该是离散化+树状数组。
# 注意把 0 也加到离散化中。

# TODO
# Run time Error
from typing import List, Tuple


def maxMatrix(n: int, m: int, queries: List[Tuple[int, int, int]]) -> List[int]:
    nums1, nums2 = [0] * n, [0] * m
    sum1, sum2 = BIT1(int(1e8 + 10)), BIT1(int(1e8 + 10))
    count1, count2 = BIT1(int(1e8 + 10)), BIT1(int(1e8 + 10))
    count1.add(0, n)
    count2.add(0, m)
    cur = 0
    res = []
    for op, qi, qv in queries:
        curNums, curCount, curSum = (nums1, count1, sum1) if op == 2 else (nums2, count2, sum2)
        otherCount, otherSum = (count2, sum2) if op == 2 else (count1, sum1)

        # 移除旧值
        preNum = curNums[qi]
        # (as min/max)
        cur -= otherCount.query(preNum) * preNum + otherSum.queryRange(preNum + 1, int(1e8 + 10))
        cur += otherCount.query(qv) * qv + otherSum.queryRange(qv + 1, int(1e8 + 10))
        curCount.add(preNum, -1)
        curSum.add(preNum, -preNum)
        # 加入新值
        curNums[qi] = qv
        curCount.add(qv, 1)
        curSum.add(qv, qv)

        res.append(cur)

    return res


if __name__ == "__main__":

    import sys

    class BIT1:
        """单点修改"""

        __slots__ = "size", "bit", "tree"

        def __init__(self, n: int):
            self.size = n
            self.tree = dict()

        def add(self, index: int, delta: int) -> None:
            # assert index >= 1, 'index must be greater than 0'
            index += 1
            while index <= self.size:
                self.tree[index] = self.tree.get(index, 0) + delta
                index += index & -index

        def query(self, index: int) -> int:
            index += 1
            if index > self.size:
                index = self.size
            res = 0
            while index > 0:
                res += self.tree.get(index, 0)
                index -= index & -index
            return res

        def queryRange(self, left: int, right: int) -> int:
            return self.query(right) - self.query(left - 1)

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n, m, q = map(int, input().split())
    Q = []
    for _ in range(q):
        op, qi, qv = map(int, input().split())
        Q.append((op, qi - 1, qv))
    print(*maxMatrix(n, m, Q), sep="\n")
