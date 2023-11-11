"""
现在有一个 n*n 的网格, 每个网格上都有一个数字, 
现在问从 (0,0) 走到 (n-1,n-1) , 并且路径上的数字异或为0的所有的方案有多少?

n<=20
折半枚举/折半搜索/半分全列挙
!以对角线为界(x + y == n -1)
C(38,19) => 2*2^n (二项式系数和)
"""


from collections import defaultdict
import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n = int(input())
    grid = [tuple(map(int, input().split())) for _ in range(n)]

    def dfs1(row: int, col: int, xor_: int) -> None:
        if row + col == n - 1:
            counter[(row, col, xor_)] += 1
            return
        if row + 1 < n:
            dfs1(row + 1, col, xor_ ^ grid[row + 1][col])
        if col + 1 < n:
            dfs1(row, col + 1, xor_ ^ grid[row][col + 1])

    def dfs2(row: int, col: int, xor_: int) -> None:
        global res
        if row + col == n - 1:
            res += counter[((row, col, xor_ ^ grid[row][col]))]
            return
        if row - 1 >= 0:
            dfs2(row - 1, col, xor_ ^ grid[row - 1][col])
        if col - 1 >= 0:
            dfs2(row, col - 1, xor_ ^ grid[row][col - 1])

    counter = defaultdict(int)
    res = 0
    dfs1(0, 0, grid[0][0])
    dfs2(n - 1, n - 1, grid[n - 1][n - 1])
    print(res)
