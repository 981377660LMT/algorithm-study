from itertools import permutations
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# H 行
# W 列の
# 2 つのグリッド A, B が与えられます。

# 1≤i≤H と
# 1≤j≤W を満たす各整数の組
# (i,j) について、
# i 行目
# j 列目にあるマスをマス
# (i,j) と呼ぶとき、 グリッド A の マス
# (i,j) には整数
# A
# i,j
# ​
#   が、 グリッド B の マス
# (i,j) には整数
# B
# i,j
# ​
#   がそれぞれ書かれています。

# あなたは「下記の
# 2 つのうちのどちらか
# 1 つを行う」という操作を好きな回数（
# 0 回でもよい）だけ繰り返します。

# 1≤i≤H−1 を満たす整数
# i を選び、グリッド A の
# i 行目と
# (i+1) 行目を入れ替える。
# 1≤i≤W−1 を満たす整数
# i を選び、グリッド A の
# i 列目と
# (i+1) 列目を入れ替える。
# 上記の操作の繰り返しによって、グリッド A をグリッド B に一致させることが可能かどうかを判定してください。 さらに、一致させることが可能な場合は、そのために行う操作回数の最小値を出力してください。

# ここで、グリッド A とグリッド B が一致しているとは、
# 1≤i≤H と
# 1≤j≤W を満たす全ての整数の組
# (i,j) について、 グリッド A の マス
# (i,j) とグリッド B の マス
# (i,j) に書かれた整数が等しいこととします。
if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    A = [list(map(int, input().split())) for _ in range(ROW)]
    B = [list(map(int, input().split())) for _ in range(ROW)]

    def minAdjacentSwap1(nums1: List[int], nums2: List[int]) -> int:
        res = 0
        for num in nums1:
            index = nums2.index(num)
            res += index
            nums2.pop(index)
        return res

    def check(a: List[List[int]], b: List[List[int]]) -> bool:
        for i in range(ROW):
            for j in range(COL):
                if a[i][j] != b[i][j]:
                    return False
        return True

    res = INF
    row1 = list(range(ROW))
    col1 = list(range(COL))
    for rowPerm in permutations(range(ROW)):
        for colPerm in permutations(range(COL)):
            a = [[0] * COL for _ in range(ROW)]
            for i in range(ROW):
                for j in range(COL):
                    a[i][j] = A[rowPerm[i]][colPerm[j]]
            if check(a, B):
                swap1 = minAdjacentSwap1(row1, list(rowPerm))
                swap2 = minAdjacentSwap1(col1, list(colPerm))
                res = min(res, swap1 + swap2)

    print(res if res != INF else -1)
