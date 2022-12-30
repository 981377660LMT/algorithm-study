from itertools import product
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# (1,2,…,N) の順列 P=(P
# 1
# ​
#  ,P
# 2
# ​
#  ,…,P
# N
# ​
#  ) が与えられます。

# すべての i (1≤i≤N) に対して、以下の値を求めてください。

# D
# i
# ​
#  =
# j
# 
# =i
# min
# ​
#  (∣P
# i
# ​
#  −P
# j
# ​
#  ∣+∣i−j∣)


from bisect import bisect_left, bisect_right
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class BIT2D:
    """二维矩形计数 更新和查询时间复杂度O(logk)

    https://www.dounaite.com/article/62af78a4b80f116a57952d98.html
    """

    def __init__(self, row: int):
        self._row = row
        self._tree = [[] for _ in range(row)]

    def add(self, row: int, col: int) -> None:
        """加入点(row,col) 注意加入过程中需要保证col递增 时间复杂度log(k)"""
        if row <= 0:
            raise ValueError("row 必须是正整数")
        while row < self._row:
            self._tree[row].append(col)
            row += row & -row

    def query(self, r1: int, c1: int, r2: int, c2: int) -> int:
        """计算矩形内的点数 时间复杂度log(k)"""
        if r1 >= self._row:
            r1 = self._row - 1
        if r2 >= self._row:
            r2 = self._row - 1
        return self._query(r2, c1, c2) - self._query(r1 - 1, c1, c2)

    def _query(self, rowUpper: int, c1: int, c2: int) -> int:
        """row不超过rowUpper,col在[c1,c2]间的点数"""
        res = 0
        index = rowUpper
        while index > 0:
            res += bisect_right(self._tree[index], c2) - bisect_left(self._tree[index], c1)
            index -= index & -index
        return res


if __name__ == "__main__":

    n = int(input())
    perm = list(map(int, input().split()))
    bit = BIT2D(int(4e5 + 10))  # 最大答案是4e5
    points = [(num, i + 1) for i, num in enumerate(perm, 1)]
    for _ in range(n):
        x, y = map(int, input().split())
    points.sort(key=lambda x: x[1])  # 按y排序加入点
    for x, y in points:
        bit.add(x, y)

    def countNGT(mid: int, x: int, y: int) -> int:
        """计算到(x,y)的切比雪夫距离不超过mid的点的个数"""
        qx, qy = (x + y + 1), (x - y + 1)
        return bit.query(qx - mid, qy - mid, qx + mid, qy + mid)

    q = int(input())
    for _ in range(q):
        qx, qy, k = map(int, input().split())
        left, right = 0, int(4e5 + 10)
        while left <= right:
            mid = (left + right) // 2
            if countNGT(mid, qx, qy) < k:
                left = mid + 1
            else:
                right = mid - 1
        print(left)
