# 折半搜索(半分全列挙)
# n<=20

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# マス (i,j) にいるとき、マス (i+1,j),(i,j+1) のいずれかに移動することができます
# マス (1,1) から移動を繰り返してマス (N,N) にたどり着く方法であって、
# 通ったマス（マス (1,1),(N,N) を含む）に書かれた整数の排他的論理和(异或)が 0 となるようなものの総数を求めてください。
if __name__ == "__main__":
    n = int(input())
    grid = [tuple(map(int, input().split())) for _ in range(n)]

    @lru_cache(None)
    def dfs(row: int, col: int, xor_: int) -> int:
        if row == n - 1 and col == n - 1:
            return 1 if xor_ == 0 else 0

        res = 0
        if row + 1 < n:
            res += dfs(row + 1, col, xor_ ^ grid[row + 1][col])
        if col + 1 < n:
            res += dfs(row, col + 1, xor_ ^ grid[row][col + 1])
        return res

    res = dfs(0, 0, grid[0][0])
    dfs.cache_clear()
    print(res)
