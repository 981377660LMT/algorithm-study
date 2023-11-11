# https://zhuanlan.zhihu.com/p/524956259

# 给你两个长度为 n 的数组 A和B (n<=2e5)
# 有一个n*n的矩阵 M[i][j]=A[i]+B[j]

# 现在给出Q次询问(Q<=2e5)，每次询问给出一个矩形范围(h1 , h2, w1 , w2)，
# 求矩形范围内的所有数的gcd

# !对于gcd我们有一个常见的公式，
# !gcd(a1 , a2, a3 ,..., an) = gcd(a1, a2- a1,a3 - a2,..., an - an-1)
# !即数组的gcd等于差分数组的gcd 。

# !左上角(0,0)右下角(n-1,n-1) 的范围的gcd可表示为
# !gcd(A[0]+B[0],*(A[2]-A[1],...,A[n-1]-A[n-2]),*(B[2]-B[1],...,B[n-1]-B[n-2]))
# 要处理差分数组的区间gcd

# !求区间gcd可以用st表或者线段树。
from math import gcd
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, q = map(int, input().split())

    nums1 = list(map(int, input().split()))
    diff1 = [0] + [abs(nums1[i] - nums1[i - 1]) for i in range(1, n)]
    st1 = SparseTable(diff1, gcd)

    nums2 = list(map(int, input().split()))
    diff2 = [0] + [abs(nums2[i] - nums2[i - 1]) for i in range(1, n)]
    st2 = SparseTable(diff2, gcd)

    for _ in range(q):
        h1, h2, w1, w2 = map(int, input().split())
        h1, h2, w1, w2 = h1 - 1, h2 - 1, w1 - 1, w2 - 1
        res = nums1[h1] + nums2[w1]
        if h2 > h1:
            res = gcd(res, st1.query(h1 + 1, h2))
        if w2 > w1:
            res = gcd(res, st2.query(w1 + 1, w2))
        print(res)


if __name__ == "__main__":
    from math import ceil, floor, log2
    from typing import Callable, Generic, List, TypeVar

    T = TypeVar("T")
    Merger = Callable[[T, T], T]

    class SparseTable(Generic[T]):
        """自定义merger的ST表"""

        __slots__ = "_n", "_dp", "_merger"

        def __init__(self, arr: List[T], merger: Merger[T]):
            n, upper = len(arr), ceil(log2(len(arr))) + 1
            self._n = n
            self._merger = merger

            dp: List[List[T]] = [[0] * upper for _ in range(n)]  # type: ignore
            for i in range(n):
                dp[i][0] = arr[i]
            for j in range(1, upper):
                for i in range(n):
                    if i + (1 << (j - 1)) >= n:
                        break
                    dp[i][j] = merger(dp[i][j - 1], dp[i + (1 << (j - 1))][j - 1])
            self._dp = dp

        def query(self, left: int, right: int) -> T:
            """[left,right]区间的最大值"""
            assert 0 <= left <= right < self._n, f"{left} {right} {self._n}"
            k = floor(log2(right - left + 1))
            return self._merger(self._dp[left][k], self._dp[right - (1 << k) + 1][k])

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
